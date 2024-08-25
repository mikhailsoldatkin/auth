package converter

import (
	"time"

	modelRepo "github.com/mikhailsoldatkin/auth/internal/repository/user/redis/model"
	"github.com/mikhailsoldatkin/auth/internal/service/user/model"
	pb "github.com/mikhailsoldatkin/auth/pkg/user_v1"
)

const (
	fieldUsername  = "username"
	fieldEmail     = "email"
	fieldRole      = "role"
	fieldUpdatedAt = "updated_at"
)

// FromRepoToService converter from Redis repository User model to service User model.
func FromRepoToService(user *modelRepo.User) *model.User {
	return &model.User{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: time.Unix(0, user.CreatedAtNs),
		UpdatedAt: time.Unix(0, user.UpdatedAtNs),
	}
}

// FromServiceToRepo converter from service User model to Redis repository User model.
func FromServiceToRepo(user *model.User) *modelRepo.User {
	now := time.Now().UnixNano()
	repoUser := &modelRepo.User{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		Role:        user.Role,
		CreatedAtNs: now,
		UpdatedAtNs: now,
	}

	return repoUser
}

// FromServiceToRepoUpdate converts a service User model to a Redis update map.
func FromServiceToRepoUpdate(updates *model.User) map[string]any {
	updateFields := make(map[string]any)

	if updates.Username != "" {
		updateFields[fieldUsername] = updates.Username
	}
	if updates.Email != "" {
		updateFields[fieldEmail] = updates.Email
	}
	if updates.Role != pb.Role_UNKNOWN.String() {
		updateFields[fieldRole] = updates.Role
	}
	updateFields[fieldUpdatedAt] = time.Now().UnixNano()

	return updateFields
}
