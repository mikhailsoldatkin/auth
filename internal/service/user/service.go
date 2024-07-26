package user

import (
	"github.com/mikhailsoldatkin/auth/internal/repository"
	"github.com/mikhailsoldatkin/auth/internal/service"
)

type serv struct {
	userRepository repository.UserRepository
}

// NewService creates a new instance of the user service.
func NewService(userRepository repository.UserRepository) service.UserService {
	return &serv{userRepository: userRepository}
}
