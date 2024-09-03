package user

import (
	"context"
)

// CheckUsersExist returns an error if one of users from provided id list doesn't exist.
func (s *userService) CheckUsersExist(ctx context.Context, ids []int64) error {
	err := s.pgRepository.CheckUsersExist(ctx, ids)
	if err != nil {
		return err
	}

	return nil
}
