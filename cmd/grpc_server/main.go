package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/mikhailsoldatkin/auth/internal/config"
	"github.com/mikhailsoldatkin/auth/internal/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/mikhailsoldatkin/auth/pkg/user_v1"
)

type server struct {
	pb.UnimplementedUserV1Server
}

// Create handles the creation of a new user in the system.
func (s *server) Create(_ context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	logger.Info("create request received: %v", req)
	return &pb.CreateResponse{Id: 1}, nil
}

// Get retrieves user data by ID.
func (s *server) Get(_ context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	logger.Info("get request received: %v", req)
	now := time.Now()

	return &pb.GetResponse{
		User: &pb.User{
			Id:        req.GetId(),
			Name:      "Mikhail Soldatkin",
			Email:     "michael.soldatin@gmail.com",
			Role:      pb.Role_ADMIN,
			CreatedAt: timestamppb.New(now),
			UpdatedAt: timestamppb.New(now),
		},
	}, nil
}

// Update modifies user data.
func (s *server) Update(_ context.Context, req *pb.UpdateRequest) (*emptypb.Empty, error) {
	logger.Info("update request received: %v", req)
	return &emptypb.Empty{}, nil
}

// Delete removes a user by ID.
func (s *server) Delete(ctx context.Context, req *pb.DeleteRequest) (*emptypb.Empty, error) {
	_ = ctx
	logger.Info("delete request received: %v", req)
	return &emptypb.Empty{}, nil
}

func main() {
	cfg := config.MustLoad()

	lis, err := net.Listen("tcp", fmt.Sprintf("%v:%d", cfg.GRPC.GRPCHost, cfg.GRPC.GRPCPort))
	if err != nil {
		logger.Fatal("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	pb.RegisterUserV1Server(s, &server{})

	logger.Info("auth server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		logger.Fatal("failed to serve: %v", err)
	}
}
