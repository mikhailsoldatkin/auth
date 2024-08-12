package user

import (
	"context"
	"fmt"

	pb "github.com/mikhailsoldatkin/auth/pkg/user_v1"
)

// Update modifies an existing user's data based on the provided update request and logs the operation.
// It updates the user data in the database, then synchronizes data in cache.
func (s *serv) Update(ctx context.Context, req *pb.UpdateRequest) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		errTx := s.pgRepository.Update(ctx, req)
		if errTx != nil {
			return errTx
		}

		errTx = s.logRepository.Log(ctx, req.GetId(), fmt.Sprintf("user with ID %d updated", req.GetId()))
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return err
	}

	err = s.redisRepository.Update(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to update user %d in cache: %v", req.GetId(), err)
	}

	return nil
}
