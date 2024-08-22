package user

import (
	"context"

	"github.com/mikhailsoldatkin/platform_common/pkg/db"

	"github.com/mikhailsoldatkin/auth/internal/repository"
	"github.com/mikhailsoldatkin/auth/internal/service"
)

var _ service.UserService = (*userService)(nil)

type userService struct {
	pgRepository    repository.UserRepository
	redisRepository repository.UserRepository
	logRepository   repository.LogRepository
	txManager       db.TxManager
}

// NewUserService creates a new instance of the user service.
func NewUserService(
	pgRepository repository.UserRepository,
	redisRepository repository.UserRepository,
	logRepository repository.LogRepository,
	txManager db.TxManager,
) service.UserService {
	return &userService{
		pgRepository:    pgRepository,
		redisRepository: redisRepository,
		logRepository:   logRepository,
		txManager:       txManager,
	}
}

// No-op implementation for LogRepository
type noOpLogRepository struct{}

func (noOpLogRepository) Log(_ context.Context, _ int64, _ string) error {
	return nil
}

// No-op implementation for TxManager
type noOpTxManager struct{}

func (noOpTxManager) ReadCommitted(ctx context.Context, f db.Handler) error {
	return f(ctx)
}

// NewMockService creates a new mock instance of the user service.
func NewMockService(deps ...any) service.UserService {
	srv := userService{
		logRepository: noOpLogRepository{},
		txManager:     noOpTxManager{},
	}

	for _, v := range deps {
		switch s := v.(type) {
		case repository.UserRepository:
			srv.pgRepository = s
			srv.redisRepository = s
		}
	}

	return &srv
}
