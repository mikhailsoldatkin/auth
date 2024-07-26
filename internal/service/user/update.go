package user

import (
	"context"

	pb "github.com/mikhailsoldatkin/auth/pkg/user_v1"
)

// Update modifies an existing user's data based on the provided update request.
func (s *serv) Update(ctx context.Context, req *pb.UpdateRequest) error {
	err := s.userRepository.Update(ctx, req)
	if err != nil {
		return err
	}
	return nil
}
