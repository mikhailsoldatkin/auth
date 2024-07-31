package user

import (
	"context"
	"fmt"

	pb "github.com/mikhailsoldatkin/auth/pkg/user_v1"
)

// Update modifies an existing user's data based on the provided update request and logs the operation.
func (s *serv) Update(ctx context.Context, req *pb.UpdateRequest) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		errTx := s.userRepository.Update(ctx, req)
		if errTx != nil {
			return errTx
		}

		errTx = s.logRepository.Log(ctx, req.GetId(), fmt.Sprintf("user with ID %d updated", req.GetId()))
		if errTx != nil {
			return errTx
		}

		return nil
	})

	return err
}
