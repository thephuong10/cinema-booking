package routes

import (
	"github.com/gin-gonic/gin"
	"showtimeservice/routes/handlers"
	"showtimeservice/services"
)

func RegisterShowTimeRoute(rg *gin.RouterGroup, ss services.IShowTimeService) {
	group := rg.Group("/showtimes")

	group.GET("/:movieId", handlers.GetShowTimeByMovieId(ss))
	group.GET("/seats/:showtimeId", handlers.GetSeatsByShowTimeId(ss))
	group.POST("/", handlers.CreateShowTime(ss))
}
