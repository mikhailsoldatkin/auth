package user

import (
	"context"

	"github.com/mikhailsoldatkin/auth/internal/customerrors"
	pb "github.com/mikhailsoldatkin/auth/pkg/user_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// CheckUsersExist returns an error if one of users from provided id list doesn't exist.
func (i *Implementation) CheckUsersExist(ctx context.Context, req *pb.CheckUsersExistRequest) (*emptypb.Empty, error) {
	if len(req.GetIds()) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "no users provided for check")
	}

	err := i.userService.CheckUsersExist(ctx, req.GetIds())
	if err != nil {
		return nil, customerrors.ConvertError(err)
	}

	return &emptypb.Empty{}, nil
}
