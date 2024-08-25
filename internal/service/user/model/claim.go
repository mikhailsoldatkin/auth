package model

import "github.com/dgrijalva/jwt-go"

// ExamplePath ...
const (
	ExamplePath = "/note_v1.NoteV1/Get"
)

// UserClaims ...
type UserClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
	Role     string `json:"role"`
}
