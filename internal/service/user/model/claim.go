package model

import "github.com/dgrijalva/jwt-go"

// UserClaims ...
type UserClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
	Role     string `json:"role"`
}
