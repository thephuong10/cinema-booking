package responses

type Room struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Rows     int    `json:"rows"`
	Columns  int    `json:"columns"`
	CinemaId string `json:"cinemaId"`
	Capacity int    `json:"capacity"`
}
