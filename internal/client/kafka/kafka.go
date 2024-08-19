package kafka

import (
	"context"

	"github.com/mikhailsoldatkin/auth/internal/client/kafka/consumer"
)

// Consumer defines the interface for a Kafka consumer.
type Consumer interface {
	Consume(ctx context.Context, topicName string, handler consumer.Handler) error
	Close() error
}
