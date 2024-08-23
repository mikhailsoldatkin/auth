package auth

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a authService) Login(ctx context.Context, username, password string) (string, error) {
	return "", status.Errorf(codes.Unimplemented, "method not implemented")
}
