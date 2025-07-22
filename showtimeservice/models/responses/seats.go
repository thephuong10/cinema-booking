package responses

type Seats struct {
	MovieId   string `json:"movieId"`
	CinemaId  string `json:"cinemaId"`
	RoomId    string `json:"roomId"`
	Name      string `json:"name"`
	Row       int    `json:"row"`
	Column    int    `json:"column"`
	Available bool   `json:"available"`
}
