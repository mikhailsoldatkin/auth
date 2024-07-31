package user

import (
	"context"
)

// Delete removes a user from the system by ID.
func (s *serv) Delete(ctx context.Context, id int64) error {
	err := s.userRepository.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
