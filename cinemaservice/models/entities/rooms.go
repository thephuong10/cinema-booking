package entities

type Room struct {
	ID       string `gorm:"primaryKey"`
	Name     string
	Rows     int
	Columns  int
	CinemaId string
	Capacity int
}
