package access

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a accessService) Check(ctx context.Context, endpoint string) error {
	return status.Errorf(codes.Unimplemented, "method not implemented")
}
