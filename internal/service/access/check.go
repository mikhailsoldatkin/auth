package access

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/mikhailsoldatkin/auth/internal/utils"
	"google.golang.org/grpc/metadata"
)

const headerAuth = "Authorization"
const prefixAuth = "Bearer "

// Check verifies whether the user has the necessary permissions to access a specific endpoint.
func (a accessService) Check(ctx context.Context, endpoint string) error {
	log.Println("Check called....")

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return fmt.Errorf("metadata is not provided")
	}

	authHeader, ok := md[headerAuth]
	if !ok || len(authHeader) == 0 {
		return fmt.Errorf("authorization header is not provided")
	}

	if !strings.HasPrefix(authHeader[0], prefixAuth) {
		return fmt.Errorf("invalid authorization header format")
	}

	accessToken := strings.TrimPrefix(authHeader[0], prefixAuth)

	claims, err := utils.VerifyToken(accessToken, []byte(a.config.TokenSecretKey))
	if err != nil {
		return fmt.Errorf("access token is invalid")
	}

	roles, err := a.userRepo.GetEndpointRoles(ctx, endpoint)
	if err != nil {
		return fmt.Errorf("failed to get roles for endpoint: %w", err)
	}

	for _, role := range roles {
		if role == claims.Role {
			return nil
		}
	}

	return fmt.Errorf("access denied")
}
