package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// validator defines an interface that requires a Validate method which returns an error if validation fails.
type validator interface {
	Validate() error
}

// ValidateInterceptor is a gRPC unary interceptor that checks if the incoming request implements
// the validator interface and calls it's Validate method.
func ValidateInterceptor(
	ctx context.Context,
	req any,
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	if val, ok := req.(validator); ok {
		if err := val.Validate(); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %v", err)
		}
	}

	return handler(ctx, req)
}
