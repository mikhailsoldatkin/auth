package app

import (
	"context"
	"log"
	"time"

	"github.com/IBM/sarama"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/mikhailsoldatkin/platform_common/pkg/cache"
	"github.com/mikhailsoldatkin/platform_common/pkg/cache/redis"
	"github.com/mikhailsoldatkin/platform_common/pkg/closer"
	"github.com/mikhailsoldatkin/platform_common/pkg/db"
	"github.com/mikhailsoldatkin/platform_common/pkg/db/pg"
	"github.com/mikhailsoldatkin/platform_common/pkg/db/transaction"

	"github.com/mikhailsoldatkin/auth/internal/api/access"
	"github.com/mikhailsoldatkin/auth/internal/api/auth"
	"github.com/mikhailsoldatkin/auth/internal/api/user"
	"github.com/mikhailsoldatkin/auth/internal/client/kafka"
	kafkaConsumer "github.com/mikhailsoldatkin/auth/internal/client/kafka/consumer"
	"github.com/mikhailsoldatkin/auth/internal/config"
	"github.com/mikhailsoldatkin/auth/internal/repository"
	logRepository "github.com/mikhailsoldatkin/auth/internal/repository/log"
	pgRepository "github.com/mikhailsoldatkin/auth/internal/repository/user/pg"
	redisRepository "github.com/mikhailsoldatkin/auth/internal/repository/user/redis"
	"github.com/mikhailsoldatkin/auth/internal/service"
	accessService "github.com/mikhailsoldatkin/auth/internal/service/access"
	authService "github.com/mikhailsoldatkin/auth/internal/service/auth"
	userSaverConsumer "github.com/mikhailsoldatkin/auth/internal/service/consumer/user_create"
	userService "github.com/mikhailsoldatkin/auth/internal/service/user"
)

// serviceProvider provides access to various services and dependencies required by the application.
// It manages database connections, Redis clients, Kafka consumers, and service instances.
type serviceProvider struct {
	config *config.Config

	dbClient  db.Client
	txManager db.TxManager

	redisPool   *redigo.Pool
	redisClient cache.RedisClient

	pgRepository    repository.UserRepository
	redisRepository repository.UserRepository
	logRepository   repository.LogRepository

	userSaverConsumer service.ConsumerService

	consumer             kafka.Consumer
	consumerGroup        sarama.ConsumerGroup
	consumerGroupHandler *kafkaConsumer.GroupHandler

	userService   service.UserService
	authService   service.AuthService
	accessService service.AccessService

	userImplementation   *user.Implementation
	authImplementation   *auth.Implementation
	accessImplementation *access.Implementation
}

// newServiceProvider creates and returns a new instance of serviceProvider.
func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

// Config returns the configuration used by the serviceProvider.
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

// DBClient returns the database client used by the serviceProvider, performs a health check.
func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.Config().DB.PostgresDSN)
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

// TxManager returns the transaction manager used by the serviceProvider.
func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

// PGRepository returns the PostgreSQL repository used by the serviceProvider.
func (s *serviceProvider) PGRepository(ctx context.Context) repository.UserRepository {
	if s.pgRepository == nil {
		s.pgRepository = pgRepository.NewRepository(s.DBClient(ctx))
	}

	return s.pgRepository
}

// RedisPool returns the Redis connection pool used by the serviceProvider.
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

// RedisClient returns the Redis client used by the serviceProvider, performs a health check.
func (s *serviceProvider) RedisClient(ctx context.Context) cache.RedisClient {
	if s.redisClient == nil {
		cl := redis.NewClient(s.RedisPool(), redis.Config(s.config.Redis))

		if err := cl.Ping(ctx); err != nil {
			log.Fatalf("failed to ping Redis: %v", err)
		}

		s.redisClient = cl
	}

	return s.redisClient
}

// RedisRepository returns the Redis repository used by the serviceProvider.
func (s *serviceProvider) RedisRepository(ctx context.Context) repository.UserRepository {
	if s.redisRepository == nil {
		s.redisRepository = redisRepository.NewRepository(s.RedisClient(ctx))
	}

	return s.redisRepository
}

// LogRepository returns the log repository used by the serviceProvider.
func (s *serviceProvider) LogRepository(ctx context.Context) repository.LogRepository {
	if s.logRepository == nil {
		s.logRepository = logRepository.NewRepository(s.DBClient(ctx))
	}

	return s.logRepository
}

// UserSaverConsumer returns the user saver consumer service used by the serviceProvider.
func (s *serviceProvider) UserSaverConsumer(ctx context.Context) service.ConsumerService {
	if s.userSaverConsumer == nil {
		s.userSaverConsumer = userSaverConsumer.NewConsumerService(
			s.PGRepository(ctx),
			s.RedisRepository(ctx),
			s.Consumer(),
			s.config.KafkaConsumer,
		)
	}

	return s.userSaverConsumer
}

// Consumer returns the Kafka consumer used by the serviceProvider.
func (s *serviceProvider) Consumer() kafka.Consumer {
	if s.consumer == nil {
		s.consumer = kafkaConsumer.NewConsumer(
			s.ConsumerGroup(),
			s.ConsumerGroupHandler(),
		)
		closer.Add(s.consumer.Close)
	}

	return s.consumer
}

// ConsumerGroup returns the Kafka consumer group used by the serviceProvider.
func (s *serviceProvider) ConsumerGroup() sarama.ConsumerGroup {
	if s.consumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			s.config.KafkaConsumer.Brokers,
			s.config.KafkaConsumer.GroupID,
			s.config.KafkaConsumer.Config,
		)

		if err != nil {
			log.Fatalf("failed to create consumer group: %v", err)
		}

		s.consumerGroup = consumerGroup
	}

	return s.consumerGroup
}

// ConsumerGroupHandler returns the Kafka consumer group handler used by the serviceProvider.
func (s *serviceProvider) ConsumerGroupHandler() *kafkaConsumer.GroupHandler {
	if s.consumerGroupHandler == nil {
		s.consumerGroupHandler = kafkaConsumer.NewGroupHandler()
	}

	return s.consumerGroupHandler
}

// UserService returns the user service used by the serviceProvider.
func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userService.NewUserService(
			s.PGRepository(ctx),
			s.RedisRepository(ctx),
			s.LogRepository(ctx),
			s.TxManager(ctx),
		)
	}

	return s.userService
}

// AuthService returns the ...
func (s *serviceProvider) AuthService(ctx context.Context) service.AuthService {
	if s.authService == nil {
		s.authService = authService.NewAuthService()
	}

	return s.authService
}

// AccessService returns the ...
func (s *serviceProvider) AccessService(ctx context.Context) service.AccessService {
	if s.accessService == nil {
		s.accessService = accessService.NewAccessService()
	}

	return s.accessService
}

// UserImplementation returns the user implementation used by the serviceProvider.
func (s *serviceProvider) UserImplementation(ctx context.Context) *user.Implementation {
	if s.userImplementation == nil {
		s.userImplementation = user.NewImplementation(s.UserService(ctx))
	}

	return s.userImplementation
}

// AuthImplementation returns the auth implementation used by the serviceProvider.
func (s *serviceProvider) AuthImplementation(ctx context.Context) *auth.Implementation {
	if s.authImplementation == nil {
		s.authImplementation = auth.NewImplementation(s.AuthService(ctx))
	}

	return s.authImplementation
}

// AccessImplementation returns the access implementation used by the serviceProvider.
func (s *serviceProvider) AccessImplementation(ctx context.Context) *access.Implementation {
	if s.accessImplementation == nil {
		s.accessImplementation = access.NewImplementation(s.AccessService(ctx))
	}

	return s.accessImplementation
}
