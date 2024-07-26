package converter

import (
	"github.com/mikhailsoldatkin/auth/internal/model"
	modelRepo "github.com/mikhailsoldatkin/auth/internal/repository/user/model"
)

// ToServiceFromRepo converter from repo User model to business logic User model. TODO doc
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
