package user

import (
	"context"

	"github.com/mikhailsoldatkin/auth/internal/model"
)

// Create creates a new user in the system.
func (s *serv) Create(ctx context.Context, data *model.User) (int64, error) {
	id, err := s.userRepository.Create(ctx, data)
	if err != nil {
		return 0, err
	}
	return id, nil
}
