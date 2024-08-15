package interceptor

import (
	"context"

	"google.golang.org/grpc"
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
			return nil, err
		}
	}

	return handler(ctx, req)
}
