package requests

type CreateTicket struct {
	ShowTimeId string  `json:"showTimeId"`
	UserId     string  `json:"userId"`
	Column     int     `json:"column"`
	Row        int     `json:"row"`
	Price      float64 `json:"price"`
}
