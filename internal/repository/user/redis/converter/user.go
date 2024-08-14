package converter

import (
	"log"
	"time"

	modelRepo "github.com/mikhailsoldatkin/auth/internal/repository/user/redis/model"
	"github.com/mikhailsoldatkin/auth/internal/service/user/model"
)

// FromRepoToService converter from Redis repository User model to service User model.
func FromRepoToService(user *modelRepo.User) *model.User {
	return &model.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: time.Unix(0, user.CreatedAtNs),
		UpdatedAt: time.Unix(0, user.UpdatedAtNs),
	}
}

// FromServiceToRepo converter from service User model to Redis repository User model.
// If createNew is true, it sets both CreatedAtNs and UpdatedAtNs. Otherwise, it only sets UpdatedAtNs.
func FromServiceToRepo(user *model.User, createNew bool) *modelRepo.User {
	now := time.Now().UnixNano()

	log.Printf("user.CreatedAt = %v", user.CreatedAt)
	log.Printf("user.UpdatedAt = %v", user.UpdatedAt)

	repoUser := &modelRepo.User{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}

	if createNew {
		repoUser.CreatedAtNs = now
		repoUser.UpdatedAtNs = now
	} else {
		repoUser.CreatedAtNs = user.CreatedAt.UnixNano()
		repoUser.UpdatedAtNs = now
	}

	return repoUser
}
