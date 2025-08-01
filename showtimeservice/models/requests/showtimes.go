package requests

import "time"

type CreateShowTime struct {
	MovieID   string    `json:"movie_id"`
	CinemaID  string    `json:"cinema_id"`
	RoomID    string    `json:"room_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Price     float64   `json:"price"`
	Status    string    `json:"status"`
}
