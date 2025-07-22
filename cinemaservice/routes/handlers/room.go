package handlers

import (
	"cinemaservice/models/requests"
	"cinemaservice/models/responses"
	"cinemaservice/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetRoomById(rs services.IRoomService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Query("id")
		room, err := rs.FindById(id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		} else {
			ctx.JSON(http.StatusOK, room)
			return
		}
	}
}

func CreateRoom(rs services.IRoomService) gin.HandlerFunc {
	return func(context *gin.Context) {
		var roomReq requests.CreateRoom
		var roomRes *responses.Room
		var err error

		err = context.ShouldBindJSON(&roomReq)

		if err != nil {
			context.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		roomRes, err = rs.Create(&roomReq)

		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		context.JSON(http.StatusOK, roomRes)

		return

	}
}
