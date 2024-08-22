package user

import (
	"context"
	"fmt"

	"github.com/mikhailsoldatkin/auth/internal/service/user/model"
)

// Update modifies an existing user's data based on the provided updates and logs the operation.
// It updates the user data in the database, then synchronizes data in cache.
func (s *userServ) Update(ctx context.Context, updates *model.User) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		errTx := s.pgRepository.Update(ctx, updates)
		if errTx != nil {
			return errTx
		}

		errTx = s.logRepository.Log(ctx, updates.ID, fmt.Sprintf("user %d updated", updates.ID))
		if errTx != nil {
			return errTx
		}

		errTx = s.redisRepository.Update(ctx, updates)
		if errTx != nil {
			return fmt.Errorf("failed to update user %d in cache: %v", updates.ID, errTx)
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
