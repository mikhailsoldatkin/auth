package user

import (
	"context"
)

// Delete removes a user from the system by ID.
func (s *serv) Delete(ctx context.Context, id int64) error {
	return s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		return s.userRepository.Delete(ctx, id)
	})
}
