package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log"
	"showtimeservice/converters"
	"showtimeservice/models/entities"
	"showtimeservice/models/requests"
	"showtimeservice/models/responses"
	"sync"
	"time"
)

type IShowTimeService interface {
	FindByMovieId(movieId string) ([]responses.ShowTime, error)
	FindSeatsByShowTimeId(id string) ([]responses.Seats, error)
	Create(req *requests.CreateShowTime) (*responses.ShowTime, error)
}

type showTimeService struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewShowTimeService(db *gorm.DB, rdb *redis.Client) IShowTimeService {
	return &showTimeService{
		db,
		rdb,
	}
}

func (ss *showTimeService) FindByMovieId(movieId string) ([]responses.ShowTime, error) {
	var res []entities.ShowTimes

	err := ss.db.Debug().Where("movie_id= ?", movieId).Find(&res).Error

	if err != nil {
		log.Printf("Couldn't fetch data from DB")
		return nil, err
	}

	return converters.ToResponse(res), err

}

func (ss *showTimeService) FindSeatsByShowTimeId(id string) ([]responses.Seats, error) {
	cacheKey := fmt.Sprintf("seats:showtime:%s", id)
	lockKey := fmt.Sprintf("lock:%s", id)
	sessionId := uuid.New().String()

	cached, err := ss.rdb.Get(context.Background(), cacheKey).Result()

	if err == nil && cached != "" {
		fmt.Printf("Get Data from Cache")
		var seats []responses.Seats
		if err := json.Unmarshal([]byte(cached), &seats); err == nil {
			return seats, nil
		}
	}

	gotLock, _ := ss.rdb.SetNX(context.Background(), lockKey, "locking", 5*time.Second).Result()

	if gotLock {

		fmt.Printf("Acquire the lock : " + sessionId + "\n")

		defer ss.rdb.Del(context.Background(), lockKey)

		seats, err := ss.buildSeatsByShowtimeId(id)

		if err != nil {
			log.Printf("Couldn't convert data from Json: %v", err)
			return nil, err
		}

		newCached, err := json.Marshal(seats)

		if err != nil {
			log.Printf("Couldn't convert data to Json: %v", err)
			return nil, err
		}

		ss.rdb.Set(context.Background(), cacheKey, newCached, 10*time.Minute)

		return seats, nil
	} else {

		fmt.Printf("Waiting..... : " + sessionId + "\n")

		for i := 0; i < 5; i++ {
			time.Sleep(1 * time.Second)
			fmt.Printf("Retrying.... : " + sessionId + "\n")
			cached, err := ss.rdb.Get(context.Background(), cacheKey).Result()
			if err == nil && cached != "" {
				var seats []responses.Seats
				if err := json.Unmarshal([]byte(cached), &seats); err == nil {
					return seats, nil
				}
			}
		}
		return nil, errors.New("system busy, Please try again")
	}

}

func (ss *showTimeService) buildSeatsByShowtimeId(id string) ([]responses.Seats, error) {
	fmt.Printf("Get Data from DB")
	var st entities.ShowTimes
	err := ss.db.Where("id = ?", id).First(&st).Error
	if err != nil {
		log.Printf("Couldn't fetch data from DB: %v", err)
		return nil, err
	}

	var (
		room       responses.Rooms
		tickets    []responses.Tickets
		err1, err2 error
	)

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		room, err1 = fetchRoomInfo(st.RoomID)
	}()

	go func() {
		defer wg.Done()
		tickets, err2 = fetchBookedTickets(st.ID)
	}()

	wg.Wait()

	if err1 != nil {
		return nil, err1
	}
	if err2 != nil {
		return nil, err2
	}

	seats := buildSeats(room, tickets, st)
	return seats, nil
}

func (ss *showTimeService) Create(req *requests.CreateShowTime) (*responses.ShowTime, error) {
	st := converters.ConvertCreateShowTimeRequestToEntity(req)

	err := ss.db.Debug().Create(&st).Error

	if err != nil {
		log.Printf("Couldn't insert data to DB")
		return nil, err
	}

	return converters.ConvertShowTimeEntityToResponse(&st), err

}

func fetchRoomInfo(roomID string) (responses.Rooms, error) {
	var room responses.Rooms
	client := resty.New()
	resp, err := client.R().
		SetResult(&room).
		Get("http://127.0.0.1:8083/cinemaservice/api/room/" + roomID)

	if err != nil {
		return room, fmt.Errorf("couldn't fetch data from cinemaserivce: %v", err)
	}
	if resp.StatusCode() != 200 {
		return room, fmt.Errorf("cinemaserivce returned code %d", resp.StatusCode())
	}
	return room, nil
}

func fetchBookedTickets(showTimeID string) ([]responses.Tickets, error) {
	var tickets []responses.Tickets
	client := resty.New()
	resp, err := client.R().
		SetResult(&tickets).
		Get("http://127.0.0.1:8084/bookingservice/api/ticket/showtime/" + showTimeID)

	if err != nil {
		return nil, fmt.Errorf("couldn't fetch data from bookingservice: %v", err)
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("bookingservice returned code %d", resp.StatusCode())
	}
	return tickets, nil
}

func buildSeats(room responses.Rooms, tickets []responses.Tickets, st entities.ShowTimes) []responses.Seats {
	bookedMap := make(map[string]bool)
	for _, t := range tickets {
		key := fmt.Sprintf("%d-%d", t.Row, t.Column)
		bookedMap[key] = true
	}

	var seats []responses.Seats
	for r := 1; r <= room.Rows; r++ {
		for c := 1; c <= room.Columns; c++ {
			key := fmt.Sprintf("%d-%d", r, c)
			seats = append(seats, responses.Seats{
				Column:    c,
				Row:       r,
				Available: !bookedMap[key],
				CinemaId:  st.CinemaID,
				RoomId:    room.Id,
				MovieId:   st.MovieID,
			})
		}
	}
	return seats
}
