package routes

import (
	"github.com/gin-gonic/gin"
	"showtimeservice/services"
)

func RegisterRouters(r *gin.Engine, ss services.IShowTimeService) {

	group := r.Group("/showtimeservice/api")

	RegisterShowTimeRoute(group, ss)

}
