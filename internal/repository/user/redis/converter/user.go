package converter

import (
	"time"

	modelRepo "github.com/mikhailsoldatkin/auth/internal/repository/user/redis/model"
	"github.com/mikhailsoldatkin/auth/internal/service/user/model"
	pb "github.com/mikhailsoldatkin/auth/pkg/user_v1"
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

// ToRepoFromService converter from service User model to Redis repository User model.
func ToRepoFromService(user *model.User) *modelRepo.User {
	return &modelRepo.User{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		Role:        user.Role,
		CreatedAtNs: time.Now().UnixNano(),
		UpdatedAtNs: time.Now().UnixNano(),
	}
}

// ToRepoFromProtobuf converter from protobuf User to Redis repository User model.
func ToRepoFromProtobuf(user *pb.UpdateRequest) *modelRepo.User {
	return &modelRepo.User{
		ID:          user.Id,
		Name:        user.Name.GetValue(),
		Email:       user.Email.GetValue(),
		Role:        user.Role.String(),
		CreatedAtNs: time.Now().UnixNano(),
		UpdatedAtNs: time.Now().UnixNano(),
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
