package auth

import (
	"github.com/mikhailsoldatkin/auth/internal/config"
	"github.com/mikhailsoldatkin/auth/internal/repository"
	"github.com/mikhailsoldatkin/auth/internal/service"
)

var _ service.AuthService = (*authService)(nil)

type authService struct {
	userRepo repository.UserRepository
	config   config.Auth
}

// NewAuthService creates a new instance of the authentication service.
func NewAuthService(
	userRepo repository.UserRepository,
	config config.Auth,
) service.AuthService {
	return &authService{
		userRepo: userRepo,
		config:   config,
	}
}
