package user_create

import (
	"context"

	"github.com/mikhailsoldatkin/auth/internal/client/kafka"
	"github.com/mikhailsoldatkin/auth/internal/config"
	"github.com/mikhailsoldatkin/auth/internal/logger"
	"github.com/mikhailsoldatkin/auth/internal/repository"
	"github.com/mikhailsoldatkin/auth/internal/service"
)

var _ service.ConsumerService = (*consumerService)(nil)

// consumerService implements the ConsumerService interface for handling Kafka messages related to user operations.
type consumerService struct {
	pgRepository    repository.UserRepository
	redisRepository repository.UserRepository
	consumer        kafka.Consumer
	config          config.KafkaConsumer
}

// NewConsumerService creates a new instance of consumerService.
// It initializes the service with a user repositories, Kafka consumer, Kafka config.
func NewConsumerService(
	pgRepository repository.UserRepository,
	redisRepository repository.UserRepository,
	consumer kafka.Consumer,
	config config.KafkaConsumer,
) service.ConsumerService {
	return &consumerService{
		pgRepository:    pgRepository,
		redisRepository: redisRepository,
		consumer:        consumer,
		config:          config,
	}
}

// RunConsumer starts the Kafka consumer to process messages.
// It listens for messages and handles errors that occur during processing.
// The method returns when the context is cancelled or an error occurs.
func (s *consumerService) RunConsumer(ctx context.Context) error {
	logger.Info("Kafka consumer is running")
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
func (s *consumerService) run(ctx context.Context) <-chan error {
	errChan := make(chan error)

	go func() {
		defer close(errChan)
		errChan <- s.consumer.Consume(ctx, s.config.Topic, s.UserCreateHandler)
	}()

	return errChan
}
