package services

import (
	"bookingservice/converters"
	"bookingservice/models/entities"
	"bookingservice/models/requests"
	"bookingservice/models/responses"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log"
	"strconv"
	"sync"
	"time"
)

type ITicketService interface {
	CreateTickets(tickets []requests.CreateTicket) (bool, error)
	FindAllByShowTimeId(showTimeId string) []responses.TicketResponse
}

type ticketService struct {
	db  *gorm.DB
	rdb *redis.Client
	//kf  *configs.KafkaProducer
}

func NewTicketService(db *gorm.DB, rdb *redis.Client) ITicketService {
	return &ticketService{
		db, rdb,
	}
}

func (ts *ticketService) CreateTickets(tickets []requests.CreateTicket) (bool, error) {

	ok, err := ts.isSeatsAvailable(tickets)

	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	var cacheKeys []string

	for _, ticket := range tickets {
		cacheKey := fmt.Sprintf("showtime:%s:seat:%d:%d:hold", ticket.ShowTimeId, ticket.Row, ticket.Column)
		cacheKeys = append(cacheKeys, cacheKey)
		suc, err := ts.rdb.SetNX(context.Background(), cacheKey, ticket.UserId, 5*time.Minute).Result()
		if err != nil {
			return false, err
		}
		if !suc {
			ts.rdb.Del(context.Background(), cacheKeys...)
			return false, nil
		}
	}

	go func() {
		showtimeCacheId := fmt.Sprintf("showtime:%s:seats", tickets[0].ShowTimeId)
		seats := make(map[string]string)
		for _, ticket := range tickets {
			key := fmt.Sprintf("%d:%d", ticket.Row, ticket.Column)
			seats[key] = "1"
		}
		ts.rdb.HSet(context.Background(), showtimeCacheId, seats)
		ts.rdb.Expire(context.Background(), showtimeCacheId, 5*time.Minute)
	}()

	//message, err := json.Marshal(messages.MessageWrapper[messages.HoldingTicket]{
	//	ID:       uuid.New().String(),
	//	Type:     "HoldingTicket",
	//	CreateAt: time.Now(),
	//	Payload: messages.HoldingTicket{
	//		Tickets: cacheKeys,
	//	},
	//})
	//if err == nil {
	//	err = ts.kf.SendMessage("topic_showtime", string(message))
	//	if err != nil {
	//		fmt.Errorf("kafka error: %v", err)
	//	}
	//}

	return true, nil
}

func (ts *ticketService) FindAllByShowTimeId(showTimeId string) []responses.TicketResponse {

	wg := sync.WaitGroup{}
	wg.Add(2)

	var (
		solds  []entities.Ticket
		result []entities.Ticket
		holds  map[string]string
	)

	go func() {
		err := ts.db.Debug().Where("show_time_id= ?", showTimeId).Find(&solds).Error
		if err != nil {
			log.Println("Redis Error : ", err)
		}
	}()

	go func() {
		var err error
		holds, err = ts.findAllByShowTimeIdFromCache(showTimeId)
		if err != nil {
			log.Println("Redis Error : ", err)
			holds = make(map[string]string)
		}
	}()

	wg.Wait()

	for _, ticket := range solds {
		key := fmt.Sprintf("%d:%d", ticket.Row, ticket.Column)
		if holds[key] != "1" {
			result = append(result, ticket)
		}
	}

	for key, val := range holds {
		if val == "1" {
			var row, column int
			_, _ = fmt.Sscanf(key, "%d:%d", &row, &column)
			result = append(result, entities.Ticket{
				Row:        row,
				Column:     column,
				ShowTimeId: showTimeId,
			})
		}
	}

	return converters.ConvertTicketEntityToResponse(result)
}

func (ts *ticketService) acquireLockAndRebuildCache(showtimeCacheId string, showTimeId string) (map[string]string, error) {
	lockKey := "lock:build:seats"
	acquire, _ := ts.rdb.SetNX(context.Background(), lockKey, 1, 5*time.Second).Result()
	if acquire {

		var soldTickets []entities.Ticket

		err := ts.db.Debug().Where("show_time_id= ?", showTimeId).Find(&soldTickets).Error

		if err != nil {
			log.Println("DB err ", err)
			return nil, err
		}

		result := make(map[string]string)

		for _, ticket := range soldTickets {
			key := fmt.Sprintf("%d:%d", ticket.Row, ticket.Column)
			result[key] = "1"
		}

		ts.rdb.HSet(context.Background(), showtimeCacheId, result)

		ts.rdb.Expire(context.Background(), showtimeCacheId, 5*time.Minute)

		ts.rdb.Del(context.Background(), lockKey)

		return result, nil

	} else {
		for i := 0; i < 5; i++ {
			time.Sleep(1 * time.Second)
			seats, _ := ts.rdb.HGetAll(context.Background(), showtimeCacheId).Result()
			if len(seats) > 0 {
				return seats, nil
			}
		}

		return nil, fmt.Errorf("system busy, please try again")
	}

}

func (ts *ticketService) findAllByShowTimeIdFromCache(showTimeId string) (map[string]string, error) {
	showtimeCacheId := fmt.Sprintf("showtime:%s:seats", showTimeId)
	seats, err := ts.rdb.HGetAll(context.Background(), showtimeCacheId).Result()

	if err != nil {
		log.Println("Redis err ", err)
		return nil, err
	}

	if len(seats) == 0 {
		seats, err = ts.acquireLockAndRebuildCache(showtimeCacheId, showTimeId)

		if err != nil {
			return nil, err
		}

	}

	return seats, nil
}

func (ts *ticketService) isSeatsAvailable(tickets []requests.CreateTicket) (bool, error) {

	seats, err := ts.findAllByShowTimeIdFromCache(tickets[0].ShowTimeId)

	if err != nil {
		log.Println("Redis err ", err)
		return false, err
	}

	for _, seat := range tickets {
		key := fmt.Sprintf("%d:%d", seat.Row, seat.Column)
		val := seats[key]
		if val != "" {
			num, _ := strconv.Atoi(val)
			if num == 1 {
				return false, nil
			}
		}
	}

	return true, nil

}
