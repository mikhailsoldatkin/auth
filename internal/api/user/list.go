package user

import (
	"context"

	"github.com/mikhailsoldatkin/auth/internal/service/user/converter"
	pb "github.com/mikhailsoldatkin/auth/pkg/user_v1"
)

// List lists users with pagination support using limit and offset.
func (i *Implementation) List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	usersServ, err := i.userService.List(ctx, req)
	if err != nil {
		return nil, err
	}

	users := make([]*pb.User, 0, len(usersServ))

	for _, userServ := range usersServ {
		users = append(users, converter.ToProtobufFromService(userServ))
	}

	return &pb.ListResponse{Users: users}, nil
}
