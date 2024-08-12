package app

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/mikhailsoldatkin/auth/internal/config"
	pb "github.com/mikhailsoldatkin/auth/pkg/user_v1"
	"github.com/mikhailsoldatkin/platform_common/pkg/closer"
)

// App represents the application with its dependencies and GRPC server.
type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
}

// NewApp initializes a new App instance with the given context and sets up the necessary dependencies.
func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

// Run starts the GRPC server and handles graceful shutdown.
func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	return a.runGRPCServer()
}

// initDeps initializes the dependencies required by the App.
func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initGRPCServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	_, err := config.Load()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
	reflection.Register(a.grpcServer)
	pb.RegisterUserV1Server(a.grpcServer, a.serviceProvider.UserImplementation(ctx))

	return nil
}

func (a *App) runGRPCServer() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.serviceProvider.config.GRPC.Port))
	if err != nil {
		return err
	}

	log.Printf("gRPC server is running on %d", a.serviceProvider.config.GRPC.Port)

	err = a.grpcServer.Serve(lis)
	if err != nil {
		return err
	}

	return nil
}
