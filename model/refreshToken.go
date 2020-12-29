package model

import (
	"github.com/dgrijalva/jwt-go"
)

type RSAKeyIdentifier struct {
	KID string `json:"kid"`
}

type RefreshTokenClaims struct {
	Token string `json:"token"`
	RSAKeyIdentifier
	jwt.StandardClaims
}

type UserTokenClaims struct {
	UserId string `json:"user_id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	RSAKeyIdentifier
	jwt.StandardClaims
}
