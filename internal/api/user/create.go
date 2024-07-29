package user

import (
	"context"

	"github.com/mikhailsoldatkin/auth/internal/service/user/converter"
	"github.com/mikhailsoldatkin/auth/internal/utils"
	pb "github.com/mikhailsoldatkin/auth/pkg/user_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Create handles the creation of a new user in the system.
func (i *Implementation) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	if err := utils.ValidatePassword(req.GetPassword(), req.GetPasswordConfirm()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "password validation failed: %v", err)
	}
	if !utils.ValidateEmail(req.GetEmail()) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid email format: %v", req.GetEmail())
	}

	user := &pb.User{
		Name:  req.Name,
		Email: req.Email,
		Role:  req.Role,
	}

	id, err := i.userService.Create(ctx, converter.ToServiceFromProtobuf(user))
	if err != nil {
		return nil, err
	}

	return &pb.CreateResponse{Id: id}, nil
}
