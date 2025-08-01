package main

import (
	"cinemaservice/configs"
	"cinemaservice/routes"
	"cinemaservice/services"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	r := gin.Default()

	db := configs.InitDB()

	cs := services.NewCinemaService(db)
	rs := services.NewRoomService(db)

	routes.RegisterRouters(r, cs, rs)

	if err := r.Run(":8083"); err != nil {
		log.Fatal("Failed to start the server: ", err)
	}

}
