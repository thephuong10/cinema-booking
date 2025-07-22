package routes

import (
	"cinemaservice/services"
	"github.com/gin-gonic/gin"
)

func RegisterRouters(r *gin.Engine, cs services.ICinemaService, rs services.IRoomService) {

	group := r.Group("/cinemaservice/api")

	RegisterCinemaRoute(group, cs)
	RegisterRoomRoute(group, rs)

}
