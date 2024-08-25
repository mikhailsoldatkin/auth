package auth

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetRefreshToken returns refresh token.
func (a authService) GetRefreshToken(ctx context.Context, refreshToken string) (string, error) {
	return "", status.Errorf(codes.Unimplemented, "method not implemented")
}
