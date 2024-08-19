package user_saver

import (
	"context"

	"github.com/mikhailsoldatkin/auth/internal/client/kafka"
	"github.com/mikhailsoldatkin/auth/internal/config"
	"github.com/mikhailsoldatkin/auth/internal/repository"
)

// Service implements the ConsumerService interface for handling Kafka messages related to user operations.
type Service struct {
	userRepository repository.UserRepository
	consumer       kafka.Consumer
	config         config.KafkaConsumer
}

// NewService creates a new instance of Service.
// It initializes the service with a user repository and Kafka consumer.
func NewService(
	userRepository repository.UserRepository,
	consumer kafka.Consumer,
	config config.KafkaConsumer,
) *Service {
	return &Service{
		userRepository: userRepository,
		consumer:       consumer,
		config:         config,
	}
}

// RunConsumer starts the Kafka consumer to process messages.
// It listens for messages and handles errors that occur during processing.
// The method returns when the context is cancelled or an error occurs.
func (s *Service) RunConsumer(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-s.run(ctx):
			if err != nil {
				return err
			}
		}
	}
}

// run initiates a new goroutine for consuming Kafka messages.
// It returns a channel that will receive errors encountered during consumption.
func (s *Service) run(ctx context.Context) <-chan error {
	errChan := make(chan error)

	go func() {
		defer close(errChan)
		errChan <- s.consumer.Consume(ctx, s.config.Topic, s.UserSaveHandler)
	}()

	return errChan
}
