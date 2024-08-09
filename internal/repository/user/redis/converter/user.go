package converter

import (
	"time"

	modelRepo "github.com/mikhailsoldatkin/auth/internal/repository/user/redis/model"
	"github.com/mikhailsoldatkin/auth/internal/service/user/model"
)

// ToServiceFromRepo converter from Redis repository User model to service User model.
func ToServiceFromRepo(user *modelRepo.User) *model.User {
	return &model.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: time.Unix(0, user.CreatedAtNs),
		UpdatedAt: time.Unix(0, user.UpdatedAtNs),
	}
}

// ToServiceFromRepoList converts list of Redis repository User models to list of service User models.
func ToServiceFromRepoList(users []*modelRepo.User) []*model.User {
	serviceUsers := make([]*model.User, len(users))
	for i, user := range users {
		serviceUsers[i] = ToServiceFromRepo(user)
	}
	return serviceUsers
}
