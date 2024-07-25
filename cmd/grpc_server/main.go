package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mikhailsoldatkin/auth/internal/config"
	"github.com/mikhailsoldatkin/auth/internal/logger"
	"github.com/mikhailsoldatkin/auth/internal/repository"
	userRepo "github.com/mikhailsoldatkin/auth/internal/repository/user"
	"github.com/mikhailsoldatkin/auth/internal/utils"
	pb "github.com/mikhailsoldatkin/auth/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	pb.UnimplementedUserV1Server
	userRepository repository.UserRepository
}

// Create handles the creation of a new user in the system.
func (s *server) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	if err := utils.ValidatePassword(req.GetPassword(), req.GetPasswordConfirm()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "password validation failed: %v", err)
	}
	if !utils.ValidateEmail(req.GetEmail()) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid email format: %v", req.GetEmail())
	}
	user := &pb.User{
		Name:  req.Name,
		Email: req.Email,
		Role:  req.Role,
	}
	id, err := s.userRepository.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	logger.Info("user %d created", id)
	return &pb.CreateResponse{Id: id}, nil
}

// Get retrieves user data by ID.
func (s *server) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	user, err := s.userRepository.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	logger.Info("user data retrieved %v", user)
	return &pb.GetResponse{User: user}, nil
}

// Update modifies user data.
func (s *server) Update(ctx context.Context, req *pb.UpdateRequest) (*emptypb.Empty, error) {
	err := s.userRepository.Update(ctx, req)
	if err != nil {
		return nil, err
	}
	logger.Info("user %d updated", req.GetId())
	return &emptypb.Empty{}, nil
}

// Delete removes a user by ID.
func (s *server) Delete(ctx context.Context, req *pb.DeleteRequest) (*emptypb.Empty, error) {
	err := s.userRepository.Delete(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	logger.Info("user %d deleted", req.GetId())
	return &emptypb.Empty{}, nil
}

// List lists users with pagination support using limit and offset.
func (s *server) List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	users, err := s.userRepository.List(ctx, req)
	if err != nil {
		return nil, err
	}
	logger.Info("users fetched")
	return &pb.ListResponse{Users: users}, nil
}

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

	s := grpc.NewServer()
	reflection.Register(s)

	userRepository := userRepo.NewRepository(pool)
	pb.RegisterUserV1Server(s, &server{userRepository: userRepository})

	logger.Info("%v listening at %v", cfg.AppName, lis.Addr())

	if err = s.Serve(lis); err != nil {
		logger.Fatal("failed to serve: %v", err)
	}
}
