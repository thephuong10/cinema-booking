package entities

import "time"

type Ticket struct {
	ID         string `gorm:"primaryKey"`
	ShowTimeId string
	Column     int
	Row        int
	Seat       string
	Price      float64
	BillId     string `gorm:"index"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
