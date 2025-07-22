package entities

import "time"

type Cinema struct {
	ID        string `gorm:"primaryKey"`
	Name      string
	Address   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
