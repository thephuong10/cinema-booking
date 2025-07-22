package routes

import (
	"cinemaservice/routes/handlers"
	"cinemaservice/services"
	"github.com/gin-gonic/gin"
)

func RegisterCinemaRoute(rg *gin.RouterGroup, cs services.ICinemaService) {
	group := rg.Group("/cinemas")

	group.GET("/:id", handlers.GetCinemaById(cs))
	group.POST("/", handlers.CreateCinema(cs))

}
