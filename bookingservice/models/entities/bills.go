package entities

import "time"

type Bill struct {
	ID        string `gorm:"primaryKey"`
	UserId    string
	Status    string
	Price     float64
	Tickets   []Ticket `gorm:"foreignKey:BillId;references:ID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
