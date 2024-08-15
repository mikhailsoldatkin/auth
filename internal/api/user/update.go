package user

import (
	"context"

	"github.com/mikhailsoldatkin/auth/internal/customerrors"
	"github.com/mikhailsoldatkin/auth/internal/service/user/converter"
	"github.com/mikhailsoldatkin/auth/internal/validators"
	pb "github.com/mikhailsoldatkin/auth/pkg/user_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Update modifies user data.
func (i *Implementation) Update(ctx context.Context, req *pb.UpdateRequest) (*emptypb.Empty, error) {
	if req.GetEmail() != nil {
		email := req.GetEmail().GetValue()
		if !validators.ValidateEmail(email) {
			return nil, status.Errorf(codes.InvalidArgument, "invalid email format: %v", email)
		}
	}

	err := i.userService.Update(ctx, converter.FromProtobufToServiceUpdate(req))
	if err != nil {
		return nil, customerrors.ConvertError(err)
	}

	return &emptypb.Empty{}, nil
}
