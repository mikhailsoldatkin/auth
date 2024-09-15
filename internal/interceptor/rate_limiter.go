package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	rateLimiter "github.com/mikhailsoldatkin/auth/internal/rate_limiter"
)

// RateLimiterInterceptor is a gRPC interceptor that applies rate limiting using a token bucket limiter.
// It checks whether a request is allowed based on the available tokens in the limiter.
type RateLimiterInterceptor struct {
	rateLimiter *rateLimiter.TokenBucketLimiter
}

// NewRateLimiterInterceptor creates a new RateLimiterInterceptor that uses
// the provided TokenBucketLimiter to control the flow of incoming gRPC requests.
func NewRateLimiterInterceptor(rateLimiter *rateLimiter.TokenBucketLimiter) *RateLimiterInterceptor {
	return &RateLimiterInterceptor{rateLimiter: rateLimiter}
}

// Unary is a gRPC unary interceptor that checks the rate limiter before allowing the request to be handled.
// If the rate limit is exceeded (i.e., no tokens are available), it returns a ResourceExhausted error.
// Otherwise, it passes the request to the provided handler.
func (r *RateLimiterInterceptor) Unary(
	ctx context.Context,
	req any,
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	if !r.rateLimiter.Allow() {
		return nil, status.Error(codes.ResourceExhausted, "too many requests")
	}

	return handler(ctx, req)
}
