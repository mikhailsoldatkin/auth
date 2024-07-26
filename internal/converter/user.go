package converter

import (
	"github.com/mikhailsoldatkin/auth/internal/model"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/mikhailsoldatkin/auth/pkg/user_v1"
)

// ToProtobufFromService converter from service User model to protobuf User model.
func ToProtobufFromService(user *model.User) *pb.User {
	return &pb.User{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      pb.Role(pb.Role_value[user.Role]),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}

// ToServiceFromProtobuf converter from protobuf User model to service User model.
func ToServiceFromProtobuf(user *pb.User) *model.User {
	return &model.User{
		ID:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role.String(),
		CreatedAt: user.CreatedAt.AsTime(),
		UpdatedAt: user.UpdatedAt.AsTime(),
	}
}
