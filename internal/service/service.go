package service

import (
	"context"

	"github.com/mikhailsoldatkin/auth/internal/service/user/model"
)

// UserService defines the interface for user-related business logic operations.
type UserService interface {
	Create(ctx context.Context, user *model.User) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, updates *model.User) error
	List(ctx context.Context, limit, offset int64) ([]*model.User, error)
}
