package repository

import (
	"context"

	"github.com/mikhailsoldatkin/auth/internal/repository/user/pg/filter"
	"github.com/mikhailsoldatkin/auth/internal/service/user/model"
)

// UserRepository defines the interface for user-related database operations.
type UserRepository interface {
	Create(ctx context.Context, user *model.User) (int64, error)
	Get(ctx context.Context, filter filter.UserFilter) (*model.User, error)
	GetEndpointRoles(ctx context.Context, endpoint string) ([]string, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, updates *model.User) error
	List(ctx context.Context, limit, offset int64) ([]*model.User, error)
	CheckUsersExist(ctx context.Context, ids []int64) error
}

// LogRepository defines the interface for logging database operations.
type LogRepository interface {
	Log(ctx context.Context, id int64, details string) error
}
