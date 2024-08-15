package user

import (
	"context"

	"github.com/mikhailsoldatkin/auth/internal/customerrors"
	"github.com/mikhailsoldatkin/auth/internal/service/user/converter"
	pb "github.com/mikhailsoldatkin/auth/pkg/user_v1"
)

// Create handles the creation of a new user in the system.
func (i *Implementation) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	id, err := i.userService.Create(ctx, converter.FromProtobufToServiceCreate(req))
	if err != nil {
		return nil, customerrors.ConvertError(err)
	}

	return &pb.CreateResponse{Id: id}, nil
}
