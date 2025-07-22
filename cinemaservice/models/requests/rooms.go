package requests

type CreateRoom struct {
	Name     string `json:"name"`
	Rows     int    `json:"rows"`
	Columns  int    `json:"columns"`
	CinemaId string `json:"cinemaId"`
	Capacity int    `json:"capacity"`
}
