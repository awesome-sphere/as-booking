package jwt

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const BEARER_TOKEN = "Bearer "

func checkTokenValidity(token *jwt.Token) (interface{}, error) {
	if _, isvalid := token.Method.(*jwt.SigningMethodHMAC); !isvalid {
		return nil, fmt.Errorf("Invalid token.")

	}
	return []byte(SECRET_KEY), nil
}

func parseToken(inputToken string) (bool, jwt.MapClaims) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(inputToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})
	return err == nil, claims
}

func IsValidJWT(c *gin.Context) (bool, jwt.MapClaims) {
	authHeader := c.GetHeader("Authorization")
	inputToken := authHeader[len(BEARER_TOKEN):]
	token, _ := jwt.Parse(inputToken, checkTokenValidity)
	switch token.Valid {
	case true:
		claims := token.Claims.(jwt.MapClaims)
		claims.VerifyIssuer(ISSUER, true)
		canParse, claims := parseToken(inputToken)
		return canParse, claims
	default:
		return false, nil
	}
}
