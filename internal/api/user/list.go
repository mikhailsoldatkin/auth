package user

import (
	"context"

	"github.com/mikhailsoldatkin/auth/internal/customerrors"
	"github.com/mikhailsoldatkin/auth/internal/service/user/converter"
	pb "github.com/mikhailsoldatkin/auth/pkg/user_v1"
)

// List lists users with pagination support using limit and offset.
func (i *Implementation) List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	users, err := i.userService.List(ctx, req)
	if err != nil {
		return nil, customerrors.ConvertError(err)
	}

	return &pb.ListResponse{Users: converter.ToProtobufFromServiceList(users)}, nil
}
