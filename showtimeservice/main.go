package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"showtimeservice/configs"
	"showtimeservice/routes"
	"showtimeservice/services"
)

func main() {
	r := gin.Default()

	db := configs.InitDB()

	rdb := configs.ConnectRedis()

	ss := services.NewShowTimeService(db, rdb)

	routes.RegisterRouters(r, ss)

	if err := r.Run(":8082"); err != nil {
		log.Fatal("Failed to start the server: ", err)
	}

}
