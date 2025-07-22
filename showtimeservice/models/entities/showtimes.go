package entities

import "time"

type ShowTimes struct {
	ID        string `gorm:"primaryKey"`
	MovieID   string
	CinemaID  string
	RoomID    string
	StartTime time.Time
	EndTime   time.Time
	Price     float64
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
