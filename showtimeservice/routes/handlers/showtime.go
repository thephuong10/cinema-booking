package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"showtimeservice/services"
)

func GetShowTimeByMovieId(ss services.IShowTimeService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		movieId := ctx.Query("movieId")
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
		showtimeId := ctx.Query("showtimeId")
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
