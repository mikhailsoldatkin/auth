package repository

import (
	"context"

	"github.com/mikhailsoldatkin/auth/internal/service/user/model"
	pb "github.com/mikhailsoldatkin/auth/pkg/user_v1"
)

// UserRepository defines the interface for user-related database operations.
type UserRepository interface {
	Create(ctx context.Context, user *model.User) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, req *pb.UpdateRequest) error
	List(ctx context.Context, req *pb.ListRequest) ([]*model.User, error)
}

// LogRepository defines the interface for logging database operations.
type LogRepository interface {
	Log(ctx context.Context, userID int64, details string) error
}
