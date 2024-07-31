package converter

import (
	modelRepo "github.com/mikhailsoldatkin/auth/internal/repository/user/model"
	"github.com/mikhailsoldatkin/auth/internal/service/user/model"
)

// ToServiceFromRepo converter from repository User model to service User model.
func ToServiceFromRepo(user *modelRepo.User) *model.User {
	return &model.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// ToServiceFromRepoList converts list of repository User models to list of service User models.
func ToServiceFromRepoList(users []*modelRepo.User) []*model.User {
	serviceUsers := make([]*model.User, len(users))
	for i, user := range users {
		serviceUsers[i] = ToServiceFromRepo(user)
	}
	return serviceUsers
}
