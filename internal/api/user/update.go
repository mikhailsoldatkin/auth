package user

import (
	"context"

	"github.com/mikhailsoldatkin/auth/internal/logger"
	pb "github.com/mikhailsoldatkin/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Update modifies user data.
func (i *Implementation) Update(ctx context.Context, req *pb.UpdateRequest) (*emptypb.Empty, error) {
	err := i.userService.Update(ctx, req)
	if err != nil {
		return nil, err
	}
	logger.Info("user %d updated", req.GetId())
	return &emptypb.Empty{}, nil
}
