package services

import (
	"bookingservice/configs"
	"bookingservice/models/entities"
	"bookingservice/models/messages"
	"bookingservice/models/requests"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"time"
)

type ITicketService interface {
}

type ticketService struct {
	db  *gorm.DB
	rdb *redis.Client
	kf  *configs.KafkaProducer
}

func buildSeatCacheKeys(tickets []requests.CreateTicket) []string {
	var keys []string
	for _, ticket := range tickets {
		key := fmt.Sprintf("ticket:%s:seat:%d:%d:available", ticket.ShowTimeId, ticket.Row, ticket.Column)
		keys = append(keys, key)
	}
	return keys
}

func buildSeatsLockKey(showTimeId string, seats [][2]int) string {
	lockKey := fmt.Sprintf("lock:ticket:%s:seats", showTimeId)
	for _, seat := range seats {
		lockKey += fmt.Sprintf(":%d:%d", seat[0], seat[1])
	}
	return lockKey
}

func (ts *ticketService) acquireLockAndCheckDB(
	lockKey, showTimeId string,
	missTickets []entities.Ticket, missSeats [][2]int,
) (bool, error) {
	ctx := context.Background()
	locked, err := ts.rdb.SetNX(ctx, lockKey, "locked", 5*time.Second).Result()
	if err != nil {
		return false, err
	}
	if !locked {
		return false, nil
	}
	defer ts.rdb.Del(ctx, lockKey)

	var soldTickets []entities.Ticket
	err = ts.db.Model(&entities.Ticket{}).
		Where("show_time_id= ?", showTimeId).
		Where("(row, column) IN ?", missSeats).
		Find(&soldTickets).Error
	if err != nil {
		return false, err
	}

	soldTicketMap := make(map[string]bool)
	for _, ticket := range soldTickets {
		key := fmt.Sprintf("%d-%d", ticket.Row, ticket.Column)
		soldTicketMap[key] = true
	}

	pipe := ts.rdb.Pipeline()
	hasSold := false
	for _, ticket := range missTickets {
		cacheKey := fmt.Sprintf("ticket:%s:seat:%d:%d:available", ticket.ShowTimeId, ticket.Row, ticket.Column)
		key := fmt.Sprintf("%d-%d", ticket.Row, ticket.Column)
		available := "true"
		if soldTicketMap[key] {
			available = "false"
			hasSold = true
		}
		pipe.Set(ctx, cacheKey, available, 5*time.Minute)
	}
	_, _ = pipe.Exec(ctx)

	return !hasSold, nil
}

func (ts *ticketService) retryPollSeatCache(keys []string, nRetry int, delay time.Duration) (bool, error) {
	ctx := context.Background()
	for i := 0; i < nRetry; i++ {
		time.Sleep(delay)
		statuses, err := ts.rdb.MGet(ctx, keys...).Result()
		if err != nil {
			return false, err
		}
		found := true
		for _, v := range statuses {
			if v == "false" {
				return false, nil
			}
			if v == nil {
				found = false
			}
		}
		if found {
			return true, nil
		}
	}
	return false, fmt.Errorf("system busy, please try again")
}

func detectCacheMiss(tickets []requests.CreateTicket, statuses []interface{}) ([]entities.Ticket, [][2]int) {
	var missTickets []entities.Ticket
	var missSeats [][2]int
	for i, status := range statuses {
		if status == "false" {
			return nil, nil
		}
		if status == nil {
			missTickets = append(missTickets, entities.Ticket{
				ShowTimeId: tickets[i].ShowTimeId,
				Row:        tickets[i].Row,
				Column:     tickets[i].Column,
			})
			missSeats = append(missSeats, [2]int{tickets[i].Row, tickets[i].Column})
		}
	}
	return missTickets, missSeats
}

func (ts *ticketService) isSeatsAvailable(tickets []requests.CreateTicket) (bool, error) {
	ctx := context.Background()
	keys := buildSeatCacheKeys(tickets)

	statuses, err := ts.rdb.MGet(ctx, keys...).Result()
	if err != nil {
		return false, err
	}

	missTickets, missSeats := detectCacheMiss(tickets, statuses)
	if len(missTickets) == 0 {
		return true, nil
	}

	lockKey := buildSeatsLockKey(tickets[0].ShowTimeId, missSeats)

	ok, err := ts.acquireLockAndCheckDB(lockKey, tickets[0].ShowTimeId, missTickets, missSeats)
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}
	return ts.retryPollSeatCache(keys, 5, 100*time.Millisecond)
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
		cacheKey := fmt.Sprintf("ticket:%s:seat:%d:%d", ticket.ShowTimeId, ticket.Row, ticket.Column)
		cacheKeys = append(cacheKeys, cacheKey)
		suc, err := ts.rdb.SetNX(context.Background(), cacheKey, ticket.UserId, 5*time.Minute).Result()
		if err != nil {
			return false, err
		}
		if !suc {
			for _, key := range cacheKeys {
				ts.rdb.Del(context.Background(), key)
			}
			return false, nil
		}
	}

	message, err := json.Marshal(messages.MessageWrapper[messages.HoldingTicket]{
		ID:       uuid.New().String(),
		Type:     "HoldingTicket",
		CreateAt: time.Now(),
		Payload: messages.HoldingTicket{
			Tickets: cacheKeys,
		},
	})
	if err == nil {
		err = ts.kf.SendMessage("topic_showtime", string(message))
		if err != nil {
			fmt.Errorf("kafka error: %v", err)
		}
	}

	return true, nil
}
