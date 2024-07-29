package user

import (
	"context"

	"github.com/mikhailsoldatkin/auth/internal/service/user/model"
)

// Get retrieves a user from the system by ID.
func (s *serv) Get(ctx context.Context, id int64) (*model.User, error) {
	var user *model.User
	errTx := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var err error
		user, err = s.userRepository.Get(ctx, id)
		return err
	})

	if errTx != nil {
		return nil, errTx
	}

	return user, nil
}
