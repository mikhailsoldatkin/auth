package auth

import (
	"context"

	"github.com/mikhailsoldatkin/auth/internal/customerrors"
	pb "github.com/mikhailsoldatkin/auth/pkg/auth_v1"
)

// GetAccessToken returns a new access token.
func (i *Implementation) GetAccessToken(ctx context.Context, req *pb.GetAccessTokenRequest) (*pb.GetAccessTokenResponse, error) {
	accessToken, err := i.authService.GetAccessToken(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, customerrors.ConvertError(err)
	}

	return &pb.GetAccessTokenResponse{AccessToken: accessToken}, nil
}
