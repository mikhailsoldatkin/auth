package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"

	"github.com/mikhailsoldatkin/auth/internal/service/user/model"
)

// GenerateToken ...
func GenerateToken(user model.User, secretKey []byte, duration time.Duration) (string, error) {
	claims := model.UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
		},
		Username: user.Name,
		Role:     user.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secretKey)
}

// VerifyToken ...
func VerifyToken(tokenStr string, secretKey []byte) (*model.UserClaims, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.Errorf("unexpected token signing method")
		}

		return secretKey, nil
	}

	token, err := jwt.ParseWithClaims(
		tokenStr,
		&model.UserClaims{},
		keyFunc,
	)

	if err != nil {
		return nil, errors.Errorf("invalid token: %s", err.Error())
	}

	claims, ok := token.Claims.(*model.UserClaims)
	if !ok {
		return nil, errors.Errorf("invalid token claims")
	}

	return claims, nil
}
