package service

import (
	"context"

	"github.com/mikhailsoldatkin/auth/internal/model"
	pb "github.com/mikhailsoldatkin/auth/pkg/user_v1"
)

type UserService interface {
	Create(ctx context.Context, data *model.User) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, req *pb.UpdateRequest) error
	List(ctx context.Context, req *pb.ListRequest) ([]*model.User, error)
}
