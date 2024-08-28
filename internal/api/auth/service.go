package auth

import (
	"github.com/mikhailsoldatkin/auth/internal/service"
	pb "github.com/mikhailsoldatkin/auth/pkg/auth_v1"
)

// Implementation provides methods for handling authentication related gRPC requests.
type Implementation struct {
	pb.UnimplementedAuthV1Server
	authService service.AuthService
}

// NewImplementation creates a new instance of Implementation with the given authentication service.
func NewImplementation(authService service.AuthService) *Implementation {
	return &Implementation{authService: authService}
}
