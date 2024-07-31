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
		var errTx error
		id, errTx = s.userRepository.Create(ctx, user)
		if errTx != nil {
			return errTx
		}

		errTx = s.logRepository.Log(ctx, id, fmt.Sprintf("user created with ID %d", id))
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return id, nil
}
