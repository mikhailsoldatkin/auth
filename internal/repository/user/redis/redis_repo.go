package redis

import (
	"context"
	"fmt"
	"strconv"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/mikhailsoldatkin/auth/internal/repository/user/pg/filter"
	"github.com/mikhailsoldatkin/platform_common/pkg/cache"
	"golang.org/x/crypto/bcrypt"

	"github.com/mikhailsoldatkin/auth/internal/customerrors"
	"github.com/mikhailsoldatkin/auth/internal/repository"
	"github.com/mikhailsoldatkin/auth/internal/repository/user/redis/converter"
	repoModel "github.com/mikhailsoldatkin/auth/internal/repository/user/redis/model"
	"github.com/mikhailsoldatkin/auth/internal/service/user/model"
)

var _ repository.UserRepository = (*repo)(nil)

const (
	userEntity = "user"
)

type repo struct {
	cl cache.RedisClient
}

// NewRepository creates a new instance of the Redis user repository.
func NewRepository(cl cache.RedisClient) repository.UserRepository {
	return &repo{cl: cl}
}

// Create stores a new user in Redis and returns user's ID.
func (r *repo) Create(ctx context.Context, user *model.User) (int64, error) {
	key := strconv.FormatInt(user.ID, 10)
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	user.Password = string(password)

	err = r.cl.HashSet(ctx, key, converter.FromServiceToRepo(user))
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

// Get retrieves a user from Redis by ID.
func (r *repo) Get(ctx context.Context, f filter.UserFilter) (*model.User, error) {
	if f.ID == nil {
		return nil, fmt.Errorf("failed to get user from cache, ID required")
	}

	key := strconv.FormatInt(*f.ID, 10)
	values, err := r.cl.HGetAll(ctx, key)
	if err != nil {
		return nil, err
	}

	if len(values) == 0 {
		return nil, customerrors.NewErrNotFound(userEntity, *f.ID)
	}

	var user repoModel.User
	err = redigo.ScanStruct(values, &user)
	if err != nil {
		return nil, err
	}

	return converter.FromRepoToService(&user), nil
}

// Delete removes a user from Redis by ID.
func (r *repo) Delete(ctx context.Context, id int64) error {
	key := strconv.FormatInt(id, 10)
	err := r.cl.Delete(ctx, key)
	if err != nil {
		return err
	}

	return nil
}

// Update modifies an existing user's data in Redis based on the provided user data.
func (r *repo) Update(ctx context.Context, updates *model.User) error {
	updateFields := converter.FromServiceToRepoUpdate(updates)

	key := strconv.FormatInt(updates.ID, 10)
	err := r.cl.HashSet(ctx, key, updateFields)
	if err != nil {
		return err
	}

	return nil
}

// List not implemented.
func (r *repo) List(_ context.Context, _, _ int64) ([]*model.User, error) {
	return nil, fmt.Errorf("method not implemented")
}

// GetEndpointRoles not implemented.
func (r *repo) GetEndpointRoles(_ context.Context, _ string) ([]string, error) {
	return nil, fmt.Errorf("method not implemented")
}

// CheckUsersExist not implemented.
func (r *repo) CheckUsersExist(_ context.Context, _ []int64) error {
	return fmt.Errorf("method not implemented")
}
