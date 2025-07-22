package handlers

import (
	"cinemaservice/models/requests"
	"cinemaservice/models/responses"
	"cinemaservice/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetCinemaById(cs services.ICinemaService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Query("id")
		cinema, err := cs.FindById(id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		} else {
			ctx.JSON(http.StatusOK, cinema)
			return
		}
	}
}

func CreateCinema(cs services.ICinemaService) gin.HandlerFunc {
	return func(context *gin.Context) {
		var cinemaReq requests.CreateCinema
		var cinemaRes *responses.Cinema
		var err error

		err = context.ShouldBindJSON(&cinemaReq)

		if err != nil {
			context.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		cinemaRes, err = cs.Create(&cinemaReq)

		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		context.JSON(http.StatusOK, cinemaRes)

		return

	}
}
