package converter

import (
	"time"

	"github.com/mikhailsoldatkin/auth/internal/service/user/model"
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

// ToProtobufFromServiceList converts a list of service User models to a list of protobuf User models.
func ToProtobufFromServiceList(users []*model.User) []*pb.User {
	protobufUsers := make([]*pb.User, len(users))
	for i, user := range users {
		protobufUsers[i] = ToProtobufFromService(user)
	}
	return protobufUsers
}

// ToServiceFromProtobuf converter from protobuf request to service User model.
func ToServiceFromProtobuf(req *pb.CreateRequest) *model.User {
	return &model.User{
		Name:      req.Name,
		Email:     req.Email,
		Role:      req.Role.String(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
