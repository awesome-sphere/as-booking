package jwt

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func IsValidJWT(c *gin.Context) (bool, jwt.MapClaims) {
	t := c.GetHeader("Authorization")
	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(
		t[7:],
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	if err != nil || !token.Valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return false, nil
	}
	return true, claims
}
