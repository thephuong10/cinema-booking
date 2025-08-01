package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"showtimeservice/configs"
	"showtimeservice/routes"
	"showtimeservice/services"
)

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:5500"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 3600,
	}))

	db := configs.InitDB()

	rdb := configs.ConnectRedis()

	ss := services.NewShowTimeService(db, rdb)

	routes.RegisterRouters(r, ss)

	if err := r.Run(":8082"); err != nil {
		log.Fatal("Failed to start the server: ", err)
	}

}
