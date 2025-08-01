package responses

type TicketResponse struct {
	Id         string  `json:"id"`
	ShowTimeId string  `json:"showTimeId"`
	Price      float64 `json:"price"`
	Column     int     `json:"column"`
	Row        int     `json:"row"`
}
