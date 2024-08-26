package access

import (
	"github.com/mikhailsoldatkin/auth/internal/config"
	"github.com/mikhailsoldatkin/auth/internal/repository"
	"github.com/mikhailsoldatkin/auth/internal/service"
)

var _ service.AccessService = (*accessService)(nil)

type accessService struct {
	userRepo repository.UserRepository
	config   config.Auth
}

// NewAccessService creates a new instance of the auth service.
func NewAccessService(
	userRepo repository.UserRepository,
	config config.Auth,
) service.AccessService {
	return &accessService{
		userRepo: userRepo,
		config:   config,
	}
}
