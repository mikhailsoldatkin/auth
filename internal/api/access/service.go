package access

import (
	"github.com/mikhailsoldatkin/auth/internal/service"
	pb "github.com/mikhailsoldatkin/auth/pkg/access_v1"
)

// Implementation ...
type Implementation struct {
	pb.UnimplementedAccessV1Server
	accessService service.AccessService
}

// NewImplementation creates a new instance of Implementation with the given access service.
func NewImplementation(accessService service.AccessService) *Implementation {
	return &Implementation{accessService: accessService}
}
