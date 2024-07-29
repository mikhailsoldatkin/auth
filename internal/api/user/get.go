package user

import (
	"context"

	"github.com/mikhailsoldatkin/auth/internal/logger"
	"github.com/mikhailsoldatkin/auth/internal/service/user/converter"
	pb "github.com/mikhailsoldatkin/auth/pkg/user_v1"
)

// Get retrieves user data by ID.
func (i *Implementation) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	user, err := i.userService.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	logger.Info("user data retrieved %v", user)
	return &pb.GetResponse{User: converter.ToProtobufFromService(user)}, nil
}
