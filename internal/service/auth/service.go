package auth

import (
	"github.com/mikhailsoldatkin/auth/internal/config"
	"github.com/mikhailsoldatkin/auth/internal/repository"
	"github.com/mikhailsoldatkin/auth/internal/service"
)

var _ service.AuthService = (*authService)(nil)

type authService struct {
	userPGRepo repository.UserRepository
	config     config.Auth
}

// NewAuthService creates a new instance of the authentication service.
func NewAuthService(
	userPGRepo repository.UserRepository,
	config config.Auth,
) service.AuthService {
	return &authService{
		userPGRepo: userPGRepo,
		config:     config,
	}
}
