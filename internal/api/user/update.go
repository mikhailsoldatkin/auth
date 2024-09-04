package user

import (
	"context"

	"github.com/mikhailsoldatkin/auth/internal/customerrors"
	"github.com/mikhailsoldatkin/auth/internal/service/user/converter"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/mikhailsoldatkin/auth/pkg/user_v1"
)

// Update modifies user data.
func (i *Implementation) Update(ctx context.Context, req *pb.UpdateRequest) (*emptypb.Empty, error) {
	err := i.userService.Update(ctx, converter.FromProtobufToServiceUpdate(req))
	if err != nil {
		return nil, customerrors.ConvertError(err)
	}

	return &emptypb.Empty{}, nil
}
