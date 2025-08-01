package routes

import (
	"bookingservice/services"
	"github.com/gin-gonic/gin"
)

func RegisterRouters(r *gin.Engine, ts services.ITicketService) {

	group := r.Group("/bookingservice/api")

	RegisterTicketRoute(group, ts)

}
