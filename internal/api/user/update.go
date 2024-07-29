package user

import (
	"context"

	pb "github.com/mikhailsoldatkin/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Update modifies user data.
func (i *Implementation) Update(ctx context.Context, req *pb.UpdateRequest) (*emptypb.Empty, error) {
	err := i.userService.Update(ctx, req)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
