package user_saver

import (
	"context"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	"github.com/mikhailsoldatkin/auth/internal/service/user/model"
)

// UserSaveHandler processes incoming Kafka messages.
// It unmarshals the message, creates a new user in the repository, and logs the result.
func (s *Service) UserSaveHandler(ctx context.Context, msg *sarama.ConsumerMessage) error {
	user := &model.User{}
	if err := json.Unmarshal(msg.Value, user); err != nil {
		log.Printf("error unmarshalling message: %v\n", err)
		return err
	}

	id, err := s.userRepository.Create(ctx, user)
	if err != nil {
		log.Printf("error creating user: %v\n", err)
		return err
	}

	log.Printf("User with ID %d created\n", id)

	return nil
}
