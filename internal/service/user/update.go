package user

import (
	"context"
	"fmt"

	pb "github.com/mikhailsoldatkin/auth/pkg/user_v1"
)

// Update modifies an existing user's data based on the provided update request and logs the operation.
func (s *serv) Update(ctx context.Context, req *pb.UpdateRequest) error {
	return s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		err := s.userRepository.Update(ctx, req)
		if err != nil {
			return err
		}

		logErr := s.userRepository.LogAction(ctx, req.GetId(), fmt.Sprintf("user with ID %d updated", req.GetId()))
		if logErr != nil {
			return logErr
		}

		return nil
	})
}
