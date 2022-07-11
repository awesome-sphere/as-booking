package main

import (
	"github.com/awesome-sphere/as-booking/jwt"
	"github.com/awesome-sphere/as-booking/models"
	"github.com/awesome-sphere/as-booking/service"
	"github.com/gin-gonic/gin"
)

func main() {
	// NOTE: Change to ReleaseMode when releasing the app
	gin.SetMode(gin.DebugMode)
	router := gin.Default()

	// initialze database
	jwt.InitializeJWTSettings()
	models.InitDatabase()

	router.POST("/book-seat", service.Booking)

	router.Run(":9009")
}
