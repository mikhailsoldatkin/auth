package converter

import (
	modelRepo "github.com/mikhailsoldatkin/auth/internal/repository/user/pg/model"
	"github.com/mikhailsoldatkin/auth/internal/service/user/model"
)

// FromRepoToService converter from Postgres repository User model to service User model.
func FromRepoToService(user *modelRepo.User) *model.User {
	return &model.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// FromRepoToServiceList converts list of Postgres repository User models to list of service User models.
func FromRepoToServiceList(users []*modelRepo.User) []*model.User {
	serviceUsers := make([]*model.User, len(users))
	for i, user := range users {
		serviceUsers[i] = FromRepoToService(user)
	}
	return serviceUsers
}
