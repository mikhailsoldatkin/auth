package user

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mikhailsoldatkin/auth/internal/customerrors"
	"github.com/mikhailsoldatkin/auth/internal/service/user/converter"
	"github.com/mikhailsoldatkin/auth/internal/validators"
	pb "github.com/mikhailsoldatkin/auth/pkg/user_v1"
)

// Create handles the creation of a new user in the system.
func (i *Implementation) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	if err := validators.ValidatePassword(req.GetPassword(), req.GetPasswordConfirm()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "password validation failed: %v", err)
	}

	id, err := i.userService.Create(ctx, converter.FromProtobufToServiceCreate(req))
	if err != nil {
		return nil, customerrors.ConvertError(err)
	}

	return &pb.CreateResponse{Id: id}, nil
}
