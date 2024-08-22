package consumer

import (
	"context"
	"log"

	"github.com/IBM/sarama"
	"github.com/pkg/errors"
)

// Consumer wraps a Sarama consumer group and a group handler.
type Consumer struct {
	consumerGroup        sarama.ConsumerGroup
	consumerGroupHandler *GroupHandler
}

// NewConsumer creates a new Kafka consumer instance.
func NewConsumer(
	consumerGroup sarama.ConsumerGroup,
	consumerGroupHandler *GroupHandler,
) *Consumer {
	return &Consumer{
		consumerGroup:        consumerGroup,
		consumerGroupHandler: consumerGroupHandler,
	}
}

// Consume starts consuming messages from the specified topic using the provided handler.
// It will continuously process messages until the context is cancelled or an error occurs.
func (c *Consumer) Consume(ctx context.Context, topicName string, handler Handler) error {
	c.consumerGroupHandler.msgHandler = handler

	return c.consume(ctx, topicName)
}

// Close shuts down the consumer group and releases resources.
func (c *Consumer) Close() error {
	return c.consumerGroup.Close()
}

// consume starts a loop to consume messages from the Kafka topic.
// This method handles rebalancing and errors.
func (c *Consumer) consume(ctx context.Context, topicName string) error {
	for {
		err := c.consumerGroup.Consume(ctx, []string{topicName}, c.consumerGroupHandler)
		if err != nil {
			if errors.Is(err, sarama.ErrClosedConsumerGroup) {
				return nil
			}

			return err
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}

		log.Printf("rebalancing...\n")
	}
}
