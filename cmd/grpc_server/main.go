package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/brianvoe/gofakeit"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mikhailsoldatkin/auth/internal/config"
	"github.com/mikhailsoldatkin/auth/internal/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/mikhailsoldatkin/auth/pkg/user_v1"
)

type server struct {
	pb.UnimplementedUserV1Server
	pool *pgxpool.Pool
}

// Create handles the creation of a new user in the system.
func (s *server) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	builder := sq.Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns(
			"name",
			"email",
			"role",
		).
		//Values(req.GetName(), req.GetEmail(), req.GetRole()).
		Values(gofakeit.Name(), gofakeit.Email(), "USER").
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to build query: %v", err)
	}

	var userID int
	err = s.pool.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to insert user: %v", err)
	}

	return &pb.CreateResponse{Id: int64(userID)}, nil
}

// Get retrieves user data by ID.
func (s *server) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	builder := sq.Select("*").
		From("users").
		Where(sq.Eq{"id": req.GetId()}).
		PlaceholderFormat(sq.Dollar).
		OrderBy("id ASC")

	query, args, err := builder.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to fetch user: %v", err)
	}

	var id int
	var name, email, role string
	var createdAt, updatedAt time.Time

	for rows.Next() {
		err = rows.Scan(&id, &name, &email, &role, &createdAt, &updatedAt)
		if err != nil {
			log.Fatalf("failed to scan note: %v", err)
		}
	}

	return &pb.GetResponse{
		User: &pb.User{
			Id:        int64(id),
			Name:      name,
			Email:     email,
			Role:      1,
			CreatedAt: timestamppb.New(createdAt),
			UpdatedAt: timestamppb.New(updatedAt),
		},
	}, nil

	//var user pb.User
	//err = s.pool.QueryRow(ctx, query, args...).Scan(&user.Id, &user.Name, &user.Email, &user.Role)
	//if err != nil {
	//	return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
	//}
	//
	//return &pb.GetResponse{User: &user}, nil
}

// Update modifies user data.
func (s *server) Update(_ context.Context, req *pb.UpdateRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "Update method is not implemented")
}

// Delete removes a user by ID.
func (s *server) Delete(ctx context.Context, req *pb.DeleteRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "Delete method is not implemented")
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
	pb.RegisterUserV1Server(s, &server{pool: pool})

	logger.Info("auth server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		logger.Fatal("failed to serve: %v", err)
	}
}
