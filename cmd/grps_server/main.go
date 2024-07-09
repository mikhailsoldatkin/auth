package main

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net"
	"time"

	pb "github.com/mikhailsoldatkin/auth/pkg/user_v1"
)

const grpcPort = 50051

type server struct {
	pb.UnimplementedUserV1Server
}

// Create handles the creation of a new user in the system.
func (s *server) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	log.Printf(color.GreenString("create request received: %v", req))
	return &pb.CreateResponse{Id: 1}, nil
}

// Get retrieves user data by ID.
func (s *server) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	log.Printf(color.GreenString("get request received: %v", req))
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
func (s *server) Update(ctx context.Context, req *pb.UpdateRequest) (*emptypb.Empty, error) {
	log.Printf(color.GreenString("update request received: %v", req))
	return &emptypb.Empty{}, nil
}

// Delete removes a user by ID.
func (s *server) Delete(ctx context.Context, req *pb.DeleteRequest) (*emptypb.Empty, error) {
	log.Printf(color.GreenString("delete request received: %v", req))
	return &emptypb.Empty{}, nil
}

func main() {

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", grpcPort))
	if err != nil {
		log.Fatalf(color.RedString("failed to listen: %v", err))
	}

	s := grpc.NewServer()
	reflection.Register(s)
	pb.RegisterUserV1Server(s, &server{})

	log.Printf(color.GreenString("server listening at %v", lis.Addr()))

	if err = s.Serve(lis); err != nil {
		log.Fatalf(color.RedString("failed to serve: %v", err))
	}
}
