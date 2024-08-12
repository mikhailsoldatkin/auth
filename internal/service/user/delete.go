package user

import (
	"context"
	"fmt"
)

// Delete removes a user from the system by ID. It deletes the user from both the database and the cache.
func (s *serv) Delete(ctx context.Context, id int64) error {
	err := s.pgRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	err = s.redisRepository.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete user %d from cache: %v", id, err)
	}

	return nil
}
