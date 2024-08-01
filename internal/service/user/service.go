package user

import (
	"github.com/mikhailsoldatkin/auth/internal/client/db"
	"github.com/mikhailsoldatkin/auth/internal/repository"
	"github.com/mikhailsoldatkin/auth/internal/service"
)

var _ service.UserService = (*serv)(nil)

type serv struct {
	userRepository repository.UserRepository
	logRepository  repository.LogRepository
	txManager      db.TxManager
}

// NewService creates a new instance of the user service.
func NewService(
	userRepository repository.UserRepository,
	logRepository repository.LogRepository,
	txManager db.TxManager,
) service.UserService {
	return &serv{
		userRepository: userRepository,
		logRepository:  logRepository,
		txManager:      txManager,
	}
}
