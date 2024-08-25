package auth

import (
	"context"

	"github.com/mikhailsoldatkin/auth/internal/customerrors"
	pb "github.com/mikhailsoldatkin/auth/pkg/auth_v1"
)

// Login authenticates a user with the provided username and password.
// Validates the credentials and, if successful, returns an access token.
func (i *Implementation) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	refreshToken, err := i.authService.Login(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		return nil, customerrors.ConvertError(err)
	}

	return &pb.LoginResponse{RefreshToken: refreshToken}, nil
}
