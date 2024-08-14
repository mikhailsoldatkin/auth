package user

import (
	"context"

	"github.com/mikhailsoldatkin/auth/internal/service/user/model"
)

// List retrieves a list of users from the system based on the provided limit and offset.
func (s *serv) List(ctx context.Context, limit, offset int64) ([]*model.User, error) {
	users, err := s.pgRepository.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	return users, nil
}
