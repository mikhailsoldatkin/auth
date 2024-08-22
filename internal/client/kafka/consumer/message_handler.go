package consumer

import (
	"context"
	"log"

	"github.com/IBM/sarama"
)

// Handler is a function type that processes Kafka messages.
type Handler func(ctx context.Context, msg *sarama.ConsumerMessage) error

// GroupHandler implements sarama.ConsumerGroupHandler and is used to handle Kafka messages.
type GroupHandler struct {
	msgHandler Handler
}

// NewGroupHandler creates a new GroupHandler instance.
func NewGroupHandler() *GroupHandler {
	return &GroupHandler{}
}

// Setup is called at the beginning of a new session before ConsumeClaim.
// It is used to set up any state needed for processing.
func (c *GroupHandler) Setup(sarama.ConsumerGroupSession) error {
	log.Println("Kafka consumer group session setup")
	return nil
}

// Cleanup is called at the end of the session after all ConsumeClaim calls have finished.
// It is used to clean up resources.
func (c *GroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	log.Println("Kafka consumer group session cleanup")
	return nil
}

// ConsumeClaim starts a loop to process messages from the claim.
// It is called when the consumer is assigned a new claim and handles messages until the session ends.
func (c *GroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				log.Printf("message channel was closed\n")
				return nil
			}

			log.Printf(
				"message claimed: value = %s, timestamp = %v, topic = %s\n",
				string(message.Value), message.Timestamp, message.Topic,
			)

			err := c.msgHandler(session.Context(), message)
			if err != nil {
				log.Printf("error handling message: %v\n", err)
				continue
			}

			session.MarkMessage(message, "")

		case <-session.Context().Done():
			log.Printf("session context done\n")
			return nil
		}
	}
}
