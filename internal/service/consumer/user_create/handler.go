package user_create

import (
	"context"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"

	"github.com/mikhailsoldatkin/auth/internal/service/user/model"
)

// UserCreateHandler processes incoming Kafka messages.
// It unmarshals the message, creates a new user in the repositories, and logs the result.
func (s *consumerService) UserCreateHandler(ctx context.Context, msg *sarama.ConsumerMessage) error {
	user := &model.User{}
	if err := json.Unmarshal(msg.Value, user); err != nil {
		log.Printf("error unmarshalling message: %v\n", err)
		return err
	}

	id, err := s.pgRepository.Create(ctx, user)
	if err != nil {
		log.Printf("error creating user id db: %v\n", err)
		return err
	}

	user.ID = id
	_, err = s.redisRepository.Create(ctx, user)
	if err != nil {
		// нужно ли возвращать ошибку? или можно скипнуть кэш
		log.Printf("error creating user in cache: %v\n", err)
		return err
	}

	log.Printf("user with ID %d created\n", id)

	return nil
}
