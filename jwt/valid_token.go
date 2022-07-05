package jwt

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func IsValidJWT(c *gin.Context) {
	t := c.GetHeader("Authorization")
	claims := &authenticationClaim{}

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
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "Valid Token!",
		})
		return
	}
}
