package responses

type Tickets struct {
	Id         string `json:"Id"`
	ShowTimeId string `json:"showTimeId"`
	Row        int    `json:"row"`
	Column     int    `json:"column"`
	Seat       string `json:"seat"`
}
