package routes

import (
	"cinemaservice/routes/handlers"
	"cinemaservice/services"
	"github.com/gin-gonic/gin"
)

func RegisterRoomRoute(rg *gin.RouterGroup, rs services.IRoomService) {
	group := rg.Group("/room")

	group.GET("/:id", handlers.GetRoomById(rs))
	group.POST("/", handlers.CreateRoom(rs))

}
