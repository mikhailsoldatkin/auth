package pg

import (
	"context"
	"strconv"
	"time"

	redigo "github.com/gomodule/redigo/redis"

	"github.com/mikhailsoldatkin/auth/internal/client/cache"
	"github.com/mikhailsoldatkin/auth/internal/repository"
	"github.com/mikhailsoldatkin/auth/internal/repository/user/redis/converter"
	modelRepo "github.com/mikhailsoldatkin/auth/internal/repository/user/redis/model"
	"github.com/mikhailsoldatkin/auth/internal/service/user/model"
)

type repo struct {
	cl cache.RedisClient
}

func NewRepository(cl cache.RedisClient) repository.UserRepository {
	return &repo{cl: cl}
}

func (r *repo) Create(ctx context.Context, userData *model.User) (int64, error) {
	id := int64(1)

	user := modelRepo.User{
		ID:          id,
		Name:        userData.Name,
		Email:       userData.Email,
		Role:        userData.Role,
		CreatedAtNs: time.Now().UnixNano(),
		UpdatedAtNs: time.Now().UnixNano(),
	}

	idStr := strconv.FormatInt(id, 10)
	err := r.cl.HashSet(ctx, idStr, user)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *repo) Get(ctx context.Context, id int64) (*model.User, error) {
	idStr := strconv.FormatInt(id, 10)
	values, err := r.cl.HGetAll(ctx, idStr)
	if err != nil {
		return nil, err
	}

	if len(values) == 0 {
		return nil, model.ErrorNoteNotFound
	}

	var user modelRepo.User
	err = redigo.ScanStruct(values, &user)
	if err != nil {
		return nil, err
	}

	return converter.ToServiceFromRepo(&user), nil
}
