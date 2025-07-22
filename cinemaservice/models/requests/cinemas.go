package requests

type CreateCinema struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}
