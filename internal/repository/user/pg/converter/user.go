package converter

import (
	"time"

	modelRepo "github.com/mikhailsoldatkin/auth/internal/repository/user/pg/model"
	"github.com/mikhailsoldatkin/auth/internal/service/user/model"
	pb "github.com/mikhailsoldatkin/auth/pkg/user_v1"
)

const (
	columnUsername  = "username"
	columnEmail     = "email"
	columnRole      = "role"
	columnUpdatedAt = "updated_at"
)

// FromRepoToService converter from Postgres repository User model to service User model.
func FromRepoToService(user *modelRepo.User) *model.User {
	return &model.User{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		Password:  user.Password,
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

// FromServiceToRepoUpdate converts a service User model to a Postgres update map.
func FromServiceToRepoUpdate(updates *model.User) map[string]any {
	updateFields := make(map[string]any)
	if updates.Username != "" {
		updateFields[columnUsername] = updates.Username
	}
	if updates.Email != "" {
		updateFields[columnEmail] = updates.Email
	}
	if updates.Role != pb.Role_UNKNOWN.String() {
		updateFields[columnRole] = updates.Role
	}
	updateFields[columnUpdatedAt] = time.Now()

	return updateFields
}
