package routes

import (
	"bookingservice/routes/handlers"
	"bookingservice/services"
	"github.com/gin-gonic/gin"
)

func RegisterTicketRoute(rg *gin.RouterGroup, ts services.ITicketService) {
	group := rg.Group("/ticket")

	group.GET("/showtime/:showTimeId", handlers.GetTicketByShowTimeId(ts))
	group.POST("/", handlers.CreateTicket(ts))

}
