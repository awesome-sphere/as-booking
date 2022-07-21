package main

import (
	"github.com/awesome-sphere/as-booking/db/models"
	"github.com/awesome-sphere/as-booking/jwt"
	"github.com/awesome-sphere/as-booking/kafka"
	"github.com/awesome-sphere/as-booking/service"
	"github.com/gin-gonic/gin"
)

func main() {
	// NOTE: Change to ReleaseMode when releasing the app
	gin.SetMode(gin.DebugMode)
	router := gin.Default()

	// initialze database
	kafka.InitKafkaTopic()
	jwt.InitializeJWTSettings()
	models.InitDatabase()

	router.POST("/booking/book-seat", service.BookSeat)
	router.POST("/booking/cancel-seat", service.CancelSeat)
	router.POST("/booking/buy-seat", service.BuySeat)
	router.POST("/booking/get-time-slots", service.GetAllTimeSlot)

	router.Run(":9009")
}
