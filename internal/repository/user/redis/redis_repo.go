package redis

import (
	"context"
	"strconv"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/mikhailsoldatkin/auth/internal/customerrors"
	pb "github.com/mikhailsoldatkin/auth/pkg/user_v1"

	"github.com/mikhailsoldatkin/auth/internal/client/cache"
	"github.com/mikhailsoldatkin/auth/internal/repository"
	"github.com/mikhailsoldatkin/auth/internal/repository/user/redis/converter"
	repoModel "github.com/mikhailsoldatkin/auth/internal/repository/user/redis/model"
	"github.com/mikhailsoldatkin/auth/internal/service/user/model"
)

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
	err := r.cl.HashSet(ctx, key, converter.ToRepoFromService(user))
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

	return converter.ToServiceFromRepo(&user), nil
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

// Update modifies an existing user's data in Redis based on the provided UpdateRequest.
func (r *repo) Update(ctx context.Context, req *pb.UpdateRequest) error {
	key := strconv.FormatInt(req.GetId(), 10)
	updates := converter.ToRepoFromProtobuf(req)

	err := r.cl.HashSet(ctx, key, updates)
	if err != nil {
		return err
	}

	return nil
}

// List retrieves all users from Redis based on the provided ListRequest.
func (r *repo) List(ctx context.Context, req *pb.ListRequest) ([]*model.User, error) {

	// var users []*model.User

	return nil, nil
}
