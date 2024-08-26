package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/mikhailsoldatkin/auth/internal/customerrors"
	"github.com/mikhailsoldatkin/auth/internal/service/user/model"
	"github.com/mikhailsoldatkin/auth/internal/utils"
)

// Login authenticates a user with the provided username and password.
// Validates the credentials and, if successful, returns an access token.
func (a authService) Login(ctx context.Context, username, password string) (string, error) {
	user, err := a.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return "", err
	}

	if !utils.VerifyPassword(user.Password, password) {
		return "", customerrors.NewErrInvalidPassword()
	}

	accessToken, err := utils.GenerateToken(
		model.User{
			Username: username,
			Role:     user.Role,
		},
		[]byte(a.config.TokenSecretKey),
		time.Duration(a.config.RefreshTokenExpirationMin)*time.Minute,
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate token")
	}

	return accessToken, nil
}
