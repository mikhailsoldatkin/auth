package access

import (
	"github.com/mikhailsoldatkin/auth/internal/service"
)

var _ service.AccessService = (*accessService)(nil)

type accessService struct{}

// NewAccessService creates a new instance of the auth service.
func NewAccessService() service.AccessService {
	return &accessService{}
}
