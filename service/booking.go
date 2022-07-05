package service

import (
	"github.com/awesome-sphere/as-booking/jwt"
	"github.com/gin-gonic/gin"
)

func Booking(c *gin.Context) {
	jwt.IsValidJWT(c)
}
