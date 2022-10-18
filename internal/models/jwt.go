package models

import "github.com/golang-jwt/jwt"

type JwtCustomClaims struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	jwt.StandardClaims
}
