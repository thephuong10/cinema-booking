package converters

import (
	"cinemaservice/models/entities"
	"cinemaservice/models/requests"
	"cinemaservice/models/responses"
	"github.com/google/uuid"
)

func ConvertCinemaEntityToResponse(entity *entities.Cinema, rooms []entities.Room) *responses.Cinema {

	var roomsRes []responses.Room

	for _, room := range rooms {
		roomsRes = append(roomsRes, responses.Room{
			ID:       room.ID,
			Name:     room.Name,
			Rows:     room.Rows,
			Columns:  room.Columns,
			Capacity: room.Capacity,
			CinemaId: room.CinemaId,
		})
	}

	return &responses.Cinema{
		ID:      entity.ID,
		Name:    entity.Name,
		Address: entity.Address,
		Rooms:   roomsRes,
	}
}

func ConvertCinemaRequestToEntity(req *requests.CreateCinema) *entities.Cinema {
	return &entities.Cinema{
		ID:      uuid.New().String(),
		Name:    req.Name,
		Address: req.Address,
	}
}
