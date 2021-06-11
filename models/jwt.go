package models

import "github.com/dgrijalva/jwt-go"

// JWTCustomClaims are custom claims extending default ones.
type JWTCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.StandardClaims
}
