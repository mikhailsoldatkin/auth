package interceptor

import (
	"context"
	"errors"

	"github.com/sony/gobreaker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CircuitBreakerInterceptor is a gRPC interceptor that integrates with a circuit breaker.
// It uses the circuit breaker pattern to handle failures and prevent system overloads by stopping requests to failing services.
type CircuitBreakerInterceptor struct {
	cb *gobreaker.CircuitBreaker
}

// NewCircuitBreakerInterceptor creates a new instance of CircuitBreakerInterceptor.
func NewCircuitBreakerInterceptor(cb *gobreaker.CircuitBreaker) *CircuitBreakerInterceptor {
	return &CircuitBreakerInterceptor{
		cb: cb,
	}
}

// Unary intercepts gRPC unary calls and applies the circuit breaker pattern.
// It executes the gRPC handler within the circuit breaker context and handles
// any errors that occur due to the circuit being in an open state.
func (c *CircuitBreakerInterceptor) Unary(
	ctx context.Context,
	req any,
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	res, err := c.cb.Execute(func() (any, error) {
		return handler(ctx, req)
	})

	if err != nil {
		if errors.Is(err, gobreaker.ErrOpenState) {
			return nil, status.Error(codes.Unavailable, "service unavailable")
		}

		return nil, err
	}

	return res, nil
}
