package model

import (
	"github.com/dgrijalva/jwt-go"
)

type RefreshToken struct {
	UserId  string
	RefreshTokenClaims
}

type RefreshTokenClaims struct {
	Token string `json:"token"`
	jwt.StandardClaims
}

type UserTokenClaims struct {
	Id string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	jwt.StandardClaims
}
