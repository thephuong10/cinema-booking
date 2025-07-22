package responses

type Rooms struct {
	Id       string `json:"Id"`
	CinemaId string `json:"cinemaId"`
	Name     string `json:"name"`
	Rows     int    `json:"rows"`
	Columns  int    `json:"columns"`
}
