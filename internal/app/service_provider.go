package app

import (
	"context"
	"log"

	"github.com/mikhailsoldatkin/auth/internal/client/db"
	"github.com/mikhailsoldatkin/auth/internal/client/db/pg"
	userRepository "github.com/mikhailsoldatkin/auth/internal/repository/user"
	"github.com/mikhailsoldatkin/auth/internal/service"
	userService "github.com/mikhailsoldatkin/auth/internal/service/user"

	"github.com/mikhailsoldatkin/auth/internal/api/user"
	"github.com/mikhailsoldatkin/auth/internal/closer"
	"github.com/mikhailsoldatkin/auth/internal/config"
	"github.com/mikhailsoldatkin/auth/internal/repository"
)

type serviceProvider struct {
	config             *config.Config
	dbClient           db.Client
	userRepository     repository.UserRepository
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

func (s *serviceProvider) DbClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.Config().Database.PostgresDSN)
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("db ping error: %s", err.Error())
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepository.NewRepository(s.DbClient(ctx))
	}

	return s.userRepository
}

func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userService.NewService(s.UserRepository(ctx))
	}

	return s.userService
}

func (s *serviceProvider) UserAPIImplementation(ctx context.Context) *user.Implementation {
	if s.userImplementation == nil {
		s.userImplementation = user.NewImplementation(s.UserService(ctx))
	}

	return s.userImplementation
}
