package user

import (
	"context"

	"github.com/mikhailsoldatkin/auth/internal/model"
)

// Get retrieves a user from the system by ID.
func (s *serv) Get(ctx context.Context, id int64) (*model.User, error) {
	user, err := s.userRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
