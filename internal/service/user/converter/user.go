package converter

import (
	"time"

	"github.com/mikhailsoldatkin/auth/internal/service/user/model"
	pb "github.com/mikhailsoldatkin/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// FromServiceToProtobuf converter from service User model to protobuf User model.
func FromServiceToProtobuf(user *model.User) *pb.User {
	return &pb.User{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      pb.Role(pb.Role_value[user.Role]),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}

// FromServiceToProtobufList converts a list of service User models to a list of protobuf User models.
func FromServiceToProtobufList(users []*model.User) []*pb.User {
	protobufUsers := make([]*pb.User, len(users))
	for i, user := range users {
		protobufUsers[i] = FromServiceToProtobuf(user)
	}
	return protobufUsers
}

// FromProtobufToServiceCreate converter from protobuf Create request to service User model.
func FromProtobufToServiceCreate(req *pb.CreateRequest) *model.User {
	now := time.Now()
	return &model.User{
		Name:      req.GetName(),
		Email:     req.GetEmail(),
		Role:      req.GetRole().String(),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// FromProtobufToServiceUpdate converter from protobuf Update request to service User model.
func FromProtobufToServiceUpdate(req *pb.UpdateRequest) *model.User {
	return &model.User{
		ID:        req.Id,
		Name:      req.GetName().GetValue(),
		Email:     req.GetEmail().GetValue(),
		Role:      req.GetRole().String(),
		UpdatedAt: time.Now(),
	}
}
