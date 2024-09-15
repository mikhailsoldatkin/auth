package rate_limiter

import (
	"context"
	"time"
)

// TokenBucketLimiter implements a token bucket rate limiter that controls
// the flow of events to a certain limit within a defined time period.
type TokenBucketLimiter struct {
	tokenBucketCh chan struct{}
}

// NewTokenBucketLimiter creates a new TokenBucketLimiter with a specified limit
// of tokens and a period during which the tokens will be replenished. The rate
// limiter allows up to 'limit' operations per 'period'. Tokens are replenished
// at regular intervals based on the period divided by the limit.
//
// Parameters:
//   - ctx: A context to control the lifecycle of the limiter (e.g., to stop replenishment).
//   - limit: The maximum number of tokens (i.e., allowed operations) in a period.
//   - period: The time period within which the limit applies.
//
// Returns:
//
//	A pointer to the created TokenBucketLimiter.
func NewTokenBucketLimiter(ctx context.Context, limit int, period time.Duration) *TokenBucketLimiter {
	limiter := &TokenBucketLimiter{
		tokenBucketCh: make(chan struct{}, limit),
	}

	for i := 0; i < limit; i++ {
		limiter.tokenBucketCh <- struct{}{}
	}

	replenishmentInterval := period.Nanoseconds() / int64(limit)
	go limiter.startPeriodicReplenishment(ctx, time.Duration(replenishmentInterval))

	return limiter
}

// startPeriodicReplenishment continuously replenishes tokens into the bucket
// at the calculated interval until the context is canceled.
//
// Parameters:
//   - ctx: A context to control when to stop the replenishment process.
//   - interval: The time between each token replenishment.
func (l *TokenBucketLimiter) startPeriodicReplenishment(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			l.tokenBucketCh <- struct{}{}
		}
	}
}

// Allow checks whether a token is available in the bucket. If a token is available,
// it is consumed and the function returns true. If no token is available, it returns false.
//
// Returns:
//
//	A boolean indicating whether the action is allowed (true) or should be rate-limited (false).
func (l *TokenBucketLimiter) Allow() bool {
	select {
	case <-l.tokenBucketCh:
		return true
	default:
		return false
	}
}
