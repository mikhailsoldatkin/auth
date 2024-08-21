package user

import (
	"context"
	"fmt"

	"github.com/mikhailsoldatkin/auth/internal/service/user/model"
)

// Get retrieves a user from the system by ID.
// It first attempts to fetch the user from cache; if unavailable, it fetches from database.
func (s *userServ) Get(ctx context.Context, id int64) (*model.User, error) {
	user, err := s.redisRepository.Get(ctx, id)
	if err == nil && user != nil {
		return user, nil
	}

	user, err = s.pgRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	_, err = s.redisRepository.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to cache user %d: %v", id, err)
	}

	return user, nil
}
