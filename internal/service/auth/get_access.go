package auth

import (
	"context"
	"time"

	"github.com/mikhailsoldatkin/auth/internal/customerrors"
	"github.com/mikhailsoldatkin/auth/internal/service/user/model"
	"github.com/mikhailsoldatkin/auth/internal/utils"
)

// GetAccessToken returns a new access token.
func (a authService) GetAccessToken(_ context.Context, refreshToken string) (string, error) {
	claims, err := utils.VerifyToken(refreshToken, []byte(a.config.TokenSecretKey))
	if err != nil {
		return "", customerrors.NewErrInvalidToken()
	}

	accessToken, err := utils.GenerateToken(
		model.User{
			Username: claims.Username,
			Role:     claims.Role,
		},
		[]byte(a.config.TokenSecretKey),
		time.Duration(a.config.AccessTokenExpirationMin)*time.Minute,
	)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
