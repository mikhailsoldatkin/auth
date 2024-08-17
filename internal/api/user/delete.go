package user

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/mikhailsoldatkin/auth/internal/customerrors"
	pb "github.com/mikhailsoldatkin/auth/pkg/user_v1"
)

// Delete removes a user by ID.
func (i *Implementation) Delete(ctx context.Context, req *pb.DeleteRequest) (*emptypb.Empty, error) {
	err := i.userService.Delete(ctx, req.GetId())
	if err != nil {
		return nil, customerrors.ConvertError(err)
	}

	return &emptypb.Empty{}, nil
}
