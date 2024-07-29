package user

import (
	"context"

	"github.com/mikhailsoldatkin/auth/internal/service/user/model"
	pb "github.com/mikhailsoldatkin/auth/pkg/user_v1"
)

// List retrieves a list of users from the system based on the provided request criteria.
func (s *serv) List(ctx context.Context, req *pb.ListRequest) ([]*model.User, error) {
	users, err := s.userRepository.List(ctx, req)
	if err != nil {
		return nil, err
	}
	return users, nil
}
