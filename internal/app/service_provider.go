package app

import (
	"context"
	"log"
	"time"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/mikhailsoldatkin/auth/internal/client/cache"
	"github.com/mikhailsoldatkin/auth/internal/client/cache/redis"
	logRepository "github.com/mikhailsoldatkin/auth/internal/repository/log"
	pgRepository "github.com/mikhailsoldatkin/auth/internal/repository/user/pg"
	redisRepository "github.com/mikhailsoldatkin/auth/internal/repository/user/redis"
	"github.com/mikhailsoldatkin/auth/internal/service"
	userService "github.com/mikhailsoldatkin/auth/internal/service/user"
	"github.com/mikhailsoldatkin/platform_common/pkg/db"
	"github.com/mikhailsoldatkin/platform_common/pkg/db/pg"
	"github.com/mikhailsoldatkin/platform_common/pkg/db/transaction"

	"github.com/mikhailsoldatkin/auth/internal/api/user"
	"github.com/mikhailsoldatkin/auth/internal/config"
	"github.com/mikhailsoldatkin/auth/internal/repository"
	"github.com/mikhailsoldatkin/platform_common/pkg/closer"
)

type serviceProvider struct {
	config *config.Config

	dbClient  db.Client
	txManager db.TxManager

	redisPool   *redigo.Pool
	redisClient cache.RedisClient

	pgRepository    repository.UserRepository
	redisRepository repository.UserRepository
	logRepository   repository.LogRepository

	userService        service.UserService
	userImplementation *user.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) Config() *config.Config {
	if s.config == nil {
		cfg, err := config.Load()
		if err != nil {
			log.Fatal(err)
		}
		s.config = cfg
	}

	return s.config
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.Config().Database.PostgresDSN)
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("failed to ping db: %s", err.Error())
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) PGRepository(ctx context.Context) repository.UserRepository {
	if s.pgRepository == nil {
		s.pgRepository = pgRepository.NewRepository(s.DBClient(ctx))
	}

	return s.pgRepository
}

func (s *serviceProvider) RedisPool() *redigo.Pool {
	if s.redisPool == nil {
		s.redisPool = &redigo.Pool{
			MaxIdle:     s.config.Redis.MaxIdle,
			MaxActive:   s.config.Redis.MaxActive,
			IdleTimeout: time.Duration(s.config.Redis.IdleTimeout),
			Dial: func() (redigo.Conn, error) {
				return redigo.Dial("tcp", s.config.Redis.Address)
			},
		}
		closer.Add(s.redisPool.Close)
	}

	return s.redisPool
}

func (s *serviceProvider) RedisClient(ctx context.Context) cache.RedisClient {
	if s.redisClient == nil {
		cl := redis.NewClient(s.RedisPool(), s.config.Redis)

		if err := cl.Ping(ctx); err != nil {
			log.Fatalf("failed to ping Redis: %v", err)
		}

		s.redisClient = cl
	}

	return s.redisClient
}

func (s *serviceProvider) RedisRepository(ctx context.Context) repository.UserRepository {
	if s.redisRepository == nil {
		s.redisRepository = redisRepository.NewRepository(s.RedisClient(ctx))
	}

	return s.redisRepository
}

func (s *serviceProvider) LogRepository(ctx context.Context) repository.LogRepository {
	if s.logRepository == nil {
		s.logRepository = logRepository.NewRepository(s.DBClient(ctx))
	}

	return s.logRepository
}

func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userService.NewService(
			s.PGRepository(ctx),
			s.RedisRepository(ctx),
			s.LogRepository(ctx),
			s.TxManager(ctx),
		)
	}

	return s.userService
}

func (s *serviceProvider) UserImplementation(ctx context.Context) *user.Implementation {
	if s.userImplementation == nil {
		s.userImplementation = user.NewImplementation(s.UserService(ctx))
	}

	return s.userImplementation
}
