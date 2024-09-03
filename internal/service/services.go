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
	CheckUsersExist(ctx context.Context, ids []int64) error
}

// ConsumerService defines the interface for running a Kafka consumer.
type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}

// AuthService provides methods for user authentication and token management.
type AuthService interface {
	Login(ctx context.Context, username, password string) (string, error)
	GetRefreshToken(ctx context.Context, oldRefreshToken string) (string, error)
	GetAccessToken(ctx context.Context, refreshToken string) (string, error)
}

// AccessService provides methods for checking access permissions for various endpoints.
type AccessService interface {
	Check(ctx context.Context, endpoint string) error
}
