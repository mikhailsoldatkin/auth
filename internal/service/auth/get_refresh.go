package auth

import (
	"context"
	"time"

	"github.com/mikhailsoldatkin/auth/internal/customerrors"
	"github.com/mikhailsoldatkin/auth/internal/service/user/model"
	"github.com/mikhailsoldatkin/auth/internal/utils"
)

// GetRefreshToken returns refresh token.
func (a authService) GetRefreshToken(_ context.Context, oldRefreshToken string) (string, error) {
	claims, err := utils.VerifyToken(oldRefreshToken, []byte(a.config.TokenSecretKey))
	if err != nil {
		return "", customerrors.NewErrInvalidToken()
	}

	refreshToken, err := utils.GenerateToken(
		model.User{
			Username: claims.Username,
			Role:     claims.Role,
		},
		[]byte(a.config.TokenSecretKey),
		time.Duration(a.config.RefreshTokenExpirationMin)*time.Minute,
	)
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}
