package repository

import (
	"context"

	pb "github.com/mikhailsoldatkin/auth/pkg/user_v1"
)

type UserRepository interface {
	Create(ctx context.Context, data *pb.User) (int64, error)
	Get(ctx context.Context, id int64) (*pb.User, error)
	//List(ctx context.Context) ([]*pb.User, error)
	//Update(ctx context.Context, id int64) error
	//Delete(ctx context.Context, id int64) error
}
