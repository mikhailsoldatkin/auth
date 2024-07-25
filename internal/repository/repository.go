package repository

import (
	"context"

	pb "github.com/mikhailsoldatkin/auth/pkg/user_v1"
)

// UserRepository defines the interface for user-related database operations.
type UserRepository interface {
	Create(ctx context.Context, data *pb.User) (int64, error)
	Get(ctx context.Context, id int64) (*pb.User, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, req *pb.UpdateRequest) error
	List(ctx context.Context, req *pb.ListRequest) ([]*pb.User, error)
}
