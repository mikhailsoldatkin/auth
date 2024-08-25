package auth

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Login authenticates a user with the provided username and password.
// Validates the credentials and, if successful, returns an access token.
func (a authService) Login(ctx context.Context, username, password string) (string, error) {
	return "", status.Errorf(codes.Unimplemented, "method not implemented")
}
