package redis

import (
	"context"
	"errors"
	"strconv"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/mikhailsoldatkin/auth/internal/client/cache"
	"github.com/mikhailsoldatkin/auth/internal/customerrors"
	"github.com/mikhailsoldatkin/auth/internal/repository"
	"github.com/mikhailsoldatkin/auth/internal/repository/user/redis/converter"
	repoModel "github.com/mikhailsoldatkin/auth/internal/repository/user/redis/model"
	"github.com/mikhailsoldatkin/auth/internal/service/user/model"
)

const (
	userEntity      = "user"
	defaultPageSize = 10
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
	err := r.cl.HashSet(ctx, key, converter.FromServiceToRepo(user, true))
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

// Get retrieves a user from Redis by ID.
func (r *repo) Get(ctx context.Context, id int64) (*model.User, error) {
	key := strconv.FormatInt(id, 10)
	values, err := r.cl.HGetAll(ctx, key)
	if err != nil {
		return nil, err
	}

	if len(values) == 0 {
		return nil, customerrors.NewErrNotFound(userEntity, id)
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
	key := strconv.FormatInt(updates.ID, 10)
	err := r.cl.HashSet(ctx, key, converter.FromServiceToRepo(updates, false))
	if err != nil {
		return err
	}

	return nil
}

// List retrieves all users from Redis based on the provided limit and offset.
func (r *repo) List(_ context.Context, _, _ int64) ([]*model.User, error) {
	return nil, errors.New("list method not implemented")
}
