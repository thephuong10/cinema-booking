package services

import (
	"cinemaservice/converters"
	"cinemaservice/models/entities"
	"cinemaservice/models/requests"
	"cinemaservice/models/responses"
	"gorm.io/gorm"
	"log"
)

type ICinemaService interface {
	Create(req *requests.CreateCinema) (*responses.Cinema, error)
	FindById(id string) (*responses.Cinema, error)
}

type cinemaService struct {
	db *gorm.DB
}

func NewCinemaService(db *gorm.DB) ICinemaService {
	return &cinemaService{
		db,
	}
}

func (cs *cinemaService) Create(req *requests.CreateCinema) (*responses.Cinema, error) {

	cinema := converters.ConvertCinemaRequestToEntity(req)

	err := cs.db.Debug().Create(&cinema).Error

	if err != nil {
		log.Printf("Couldn't insert data to DB: %v", err)
		return nil, err
	}

	return converters.ConvertCinemaEntityToResponse(cinema, nil), nil

}

func (cs *cinemaService) FindById(id string) (*responses.Cinema, error) {
	var cinema entities.Cinema

	err := cs.db.Debug().Where("id= ?", id).First(&cinema).Error

	if err != nil {
		log.Printf("Couldn't fetch cinema data from DB: %v", err)
		return nil, err
	}

	var rooms []entities.Room

	err = cs.db.Debug().Where("cinema_id= ?", id).Find(&rooms).Error

	if err != nil {
		log.Printf("Couldn't fetch room data from DB: %v", err)
	}

	return converters.ConvertCinemaEntityToResponse(&cinema, rooms), nil
}
