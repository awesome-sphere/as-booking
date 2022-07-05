package jwt

import "github.com/dgrijalva/jwt-go"

type authenticationClaim struct {
	ID      int  `json:"id"`
	IsAdmin bool `json:"is_admin"`
	jwt.StandardClaims
}
