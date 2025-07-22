package converters

import (
	"cinemaservice/models/entities"
	"cinemaservice/models/requests"
	"cinemaservice/models/responses"
	"github.com/google/uuid"
)

func ConvertRoomEntityToResponse(entity *entities.Room) *responses.Room {
	return &responses.Room{
		ID:       entity.ID,
		Name:     entity.Name,
		Rows:     entity.Rows,
		Columns:  entity.Columns,
		Capacity: entity.Capacity,
		CinemaId: entity.CinemaId,
	}
}

func ConvertRoomRequestToEntity(req *requests.CreateRoom) *entities.Room {
	return &entities.Room{
		ID:       uuid.New().String(),
		Name:     req.Name,
		CinemaId: req.CinemaId,
		Rows:     req.Rows,
		Columns:  req.Columns,
		Capacity: req.Capacity,
	}
}
