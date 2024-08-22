package auth

import (
	"github.com/mikhailsoldatkin/auth/internal/service"
)

var _ service.AuthService = (*authService)(nil)

type authService struct{}

// NewAuthService creates a new instance of the auth service.
func NewAuthService() service.AuthService {
	return &authService{}
}
