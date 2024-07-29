package user

import (
	"context"
	"fmt"

	"github.com/mikhailsoldatkin/auth/internal/service/user/model"
)

// Create creates a new user in the system and logs the operation.
func (s *serv) Create(ctx context.Context, user *model.User) (int64, error) {
	var id int64
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var err error
		id, err = s.userRepository.Create(ctx, user)
		if err != nil {
			return err
		}

		errLog := s.userRepository.LogAction(ctx, id, fmt.Sprintf("user created with ID %d", id))
		if errLog != nil {
			return errLog
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return id, nil
}
