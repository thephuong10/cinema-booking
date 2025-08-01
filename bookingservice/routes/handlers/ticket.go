package handlers

import (
	"bookingservice/models/requests"
	"bookingservice/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateTicket(ts services.ITicketService) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var req []requests.CreateTicket

		err := ctx.ShouldBindJSON(&req)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		} else {

			suc, err := ts.CreateTickets(req)

			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			ctx.JSON(http.StatusOK, gin.H{
				"success": suc,
			})
			return

		}
	}
}

func GetTicketByShowTimeId(ts services.ITicketService) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		showTimeId := ctx.Param("showTimeId")

		tickets := ts.FindAllByShowTimeId(showTimeId)

		ctx.JSON(http.StatusOK, tickets)

		return

	}
}
