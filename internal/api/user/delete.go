package user

import (
	"context"

	pb "github.com/mikhailsoldatkin/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Delete removes a user by ID.
func (i *Implementation) Delete(ctx context.Context, req *pb.DeleteRequest) (*emptypb.Empty, error) {
	err := i.userService.Delete(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
