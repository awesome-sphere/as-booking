package service

import (
	"net/http"

	"github.com/awesome-sphere/as-booking/jwt"
	"github.com/gin-gonic/gin"
)

func Booking(c *gin.Context) {
	is_valid, claimed_token := jwt.IsValidJWT(c)
	if is_valid {
		c.JSON(http.StatusOK, gin.H{
			"message": claimed_token["user_id"],
		})
	}
}
