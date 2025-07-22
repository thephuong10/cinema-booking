package responses

type Cinema struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Rooms   []Room `json:"rooms"`
}
