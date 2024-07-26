package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/jackc/pgx/v4/pgxpool"
	userAPI "github.com/mikhailsoldatkin/auth/internal/api/user"
	"github.com/mikhailsoldatkin/auth/internal/config"
	"github.com/mikhailsoldatkin/auth/internal/logger"
	userRepository "github.com/mikhailsoldatkin/auth/internal/repository/user"
	userService "github.com/mikhailsoldatkin/auth/internal/service/user"
	pb "github.com/mikhailsoldatkin/auth/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.MustLoad()
	ctx := context.Background()

	pool, err := pgxpool.Connect(ctx, cfg.Database.PostgresDSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPC.GRPCPort))
	if err != nil {
		logger.Fatal("failed to listen: %v", err)
	}

	userRepo := userRepository.NewRepository(pool)
	userSrv := userService.NewService(userRepo)

	s := grpc.NewServer()
	reflection.Register(s)
	pb.RegisterUserV1Server(s, userAPI.NewImplementation(userSrv))

	logger.Info("%v listening at %v", cfg.AppName, lis.Addr())

	if err = s.Serve(lis); err != nil {
		logger.Fatal("failed to serve: %v", err)
	}
}
