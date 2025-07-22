package services

import (
	"cinemaservice/converters"
	"cinemaservice/models/entities"
	"cinemaservice/models/requests"
	"cinemaservice/models/responses"
	"gorm.io/gorm"
	"log"
)

type IRoomService interface {
	Create(req *requests.CreateRoom) (*responses.Room, error)
	FindById(id string) (*responses.Room, error)
}

type roomService struct {
	db *gorm.DB
}

func NewRoomService(db *gorm.DB) IRoomService {
	return &roomService{
		db,
	}
}

func (cs *roomService) Create(req *requests.CreateRoom) (*responses.Room, error) {

	room := converters.ConvertRoomRequestToEntity(req)

	err := cs.db.Debug().Create(&room).Error

	if err != nil {
		log.Printf("Couldn't insert data to DB: %v", err)
		return nil, err
	}

	return converters.ConvertRoomEntityToResponse(room), nil

}

func (cs *roomService) FindById(id string) (*responses.Room, error) {
	var room entities.Room

	err := cs.db.Debug().Where("id= ?", id).First(&room).Error

	if err != nil {
		log.Printf("Couldn't fetch room data from DB: %v", err)
		return nil, err
	}

	return converters.ConvertRoomEntityToResponse(&room), nil
}
