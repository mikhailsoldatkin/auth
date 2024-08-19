package kafka

import (
	"context"

	"github.com/mikhailsoldatkin/auth/internal/client/kafka/consumer"
)

// Consumer defines the interface for a Kafka consumer.
type Consumer interface {
	// Consume starts consuming messages from the specified topic using the provided handler.
	Consume(ctx context.Context, topicName string, handler consumer.Handler) error

	// Close shuts down the consumer and releases resources.
	Close() error
}
