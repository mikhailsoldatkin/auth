package auth

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetAccessToken returns a new access token.
func (a authService) GetAccessToken(ctx context.Context, accessToken string) (string, error) {
	return "", status.Errorf(codes.Unimplemented, "method not implemented")
}
