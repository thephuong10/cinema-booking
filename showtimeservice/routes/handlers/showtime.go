package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"showtimeservice/models/requests"
	"showtimeservice/models/responses"
	"showtimeservice/services"
)

func GetShowTimeByMovieId(ss services.IShowTimeService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		movieId := ctx.Param("movieId")
		showTimes, err := ss.FindByMovieId(movieId)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		} else {
			ctx.JSON(http.StatusOK, showTimes)
			return
		}
	}
}

func GetSeatsByShowTimeId(ss services.IShowTimeService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		showtimeId := ctx.Param("showtimeId")
		seats, err := ss.FindSeatsByShowTimeId(showtimeId)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		} else {
			ctx.JSON(http.StatusOK, seats)
			return
		}
	}
}

func CreateShowTime(ss services.IShowTimeService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req requests.CreateShowTime

		err := ctx.ShouldBindJSON(&req)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		var res *responses.ShowTime

		res, err = ss.Create(&req)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, res)
		return

	}
}
