package converter

import (
	"github.com/mikhailsoldatkin/auth/internal/repository/user/model"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/mikhailsoldatkin/auth/pkg/user_v1"
)

func ToUserFromRepo(user *model.User) *pb.User {
	return &pb.User{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      pb.Role(pb.Role_value[user.Role]),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}
