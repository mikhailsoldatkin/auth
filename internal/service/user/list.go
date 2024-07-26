package user

import (
	"context"

	"github.com/mikhailsoldatkin/auth/internal/model"
	pb "github.com/mikhailsoldatkin/auth/pkg/user_v1"
)

func (s *serv) List(ctx context.Context, req *pb.ListRequest) ([]*model.User, error) {
	users, err := s.userRepository.List(ctx, req)
	if err != nil {
		return nil, err
	}
	return users, nil
}
