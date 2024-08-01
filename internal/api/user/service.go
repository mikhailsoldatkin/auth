package user

import (
	"github.com/mikhailsoldatkin/auth/internal/service"
	pb "github.com/mikhailsoldatkin/auth/pkg/user_v1"
)

// Implementation provides methods for handling user-related gRPC requests.
type Implementation struct {
	pb.UnimplementedUserV1Server
	userService service.UserService
}

// NewImplementation creates a new instance of Implementation with the given user service.
func NewImplementation(userService service.UserService) *Implementation {
	return &Implementation{userService: userService}
}
