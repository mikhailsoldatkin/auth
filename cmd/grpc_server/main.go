package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"regexp"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mikhailsoldatkin/auth/internal/config"
	"github.com/mikhailsoldatkin/auth/internal/logger"
	"github.com/mikhailsoldatkin/auth/internal/repository"
	userRepo "github.com/mikhailsoldatkin/auth/internal/repository/user"
	pb "github.com/mikhailsoldatkin/auth/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

const (
	defaultPageSize   = 10
	emailRegex        = `^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`
	passwordMinLength = 8
)

type server struct {
	pb.UnimplementedUserV1Server
	userRepository repository.UserRepository
}

// validateEmail checks if the given email address is in valid format.
func validateEmail(email string) bool {
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

// validatePassword provides simple password validation.
func validatePassword(password, passwordConfirm string) error {
	if len(password) < passwordMinLength {
		return errors.New("password must be at least 8 characters long")
	}
	if password != passwordConfirm {
		return errors.New("passwords don't match")
	}
	return nil
}

// Create handles the creation of a new user in the system.
func (s *server) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	if err := validatePassword(req.GetPassword(), req.GetPasswordConfirm()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "password validation failed: %v", err)
	}
	if !validateEmail(req.GetEmail()) {
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

//// Update modifies user data.
//func (s *server) Update(ctx context.Context, req *pb.UpdateRequest) (*emptypb.Empty, error) {
//	if err := s.checkUserExists(ctx, req.GetId()); err != nil {
//		return nil, err
//	}
//
//	updateFields := make(map[string]any)
//
//	if req.GetName() != nil {
//		updateFields[columnName] = req.GetName().GetValue()
//	}
//	if req.GetEmail() != nil {
//		email := req.GetEmail().GetValue()
//		if !validateEmail(email) {
//			return nil, status.Errorf(codes.InvalidArgument, "invalid email format: %v", email)
//		}
//		updateFields[columnEmail] = email
//	}
//	if req.GetRole().String() != "" {
//		updateFields[columnRole] = req.GetRole().String()
//	}
//
//	if len(updateFields) == 0 {
//		return nil, status.Errorf(codes.InvalidArgument, "no fields to update")
//	}
//
//	updateFields[columnUpdatedAt] = time.Now()
//
//	builder := sq.Update(tableUsers).
//		SetMap(updateFields).
//		Where(sq.Eq{columnID: req.GetId()}).
//		PlaceholderFormat(sq.Dollar)
//
//	query, args, err := builder.ToSql()
//	if err != nil {
//		return nil, status.Errorf(codes.Internal, "failed to build update query: %v", err)
//	}
//
//	_, err = s.userRepo.Exec(ctx, query, args...)
//	if err != nil {
//		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
//	}
//
//	logger.Info("user %d updated", req.GetId())
//
//	return &emptypb.Empty{}, nil
//}
//
//// Delete removes a user by ID.
//func (s *server) Delete(ctx context.Context, req *pb.DeleteRequest) (*emptypb.Empty, error) {
//	builder := sq.Delete(tableUsers).
//		Where(sq.Eq{columnID: req.GetId()}).
//		PlaceholderFormat(sq.Dollar)
//
//	query, args, err := builder.ToSql()
//	if err != nil {
//		return nil, status.Errorf(codes.Internal, "failed to build delete query: %v", err)
//	}
//
//	_, err = s.userRepo.Exec(ctx, query, args...)
//	if err != nil {
//		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
//	}
//
//	logger.Info("user %d deleted", req.GetId())
//
//	return &emptypb.Empty{}, nil
//}
//
//// List lists users with pagination support using limit and offset.
//func (s *server) List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
//	limit := int(req.GetLimit())
//	fmt.Println(limit)
//	if limit <= 0 {
//		limit = defaultPageSize
//	}
//
//	builder := sq.Select("*").
//		From(tableUsers).
//		OrderBy(columnID).
//		Limit(uint64(limit)).
//		Offset(uint64(int(req.GetOffset()))).
//		PlaceholderFormat(sq.Dollar)
//
//	query, args, err := builder.ToSql()
//	if err != nil {
//		return nil, status.Errorf(codes.Internal, "failed to build query: %v", err)
//	}
//
//	rows, err := s.userRepo.Query(ctx, query, args...)
//	if err != nil {
//		return nil, status.Errorf(codes.Internal, "failed to list users: %v", err)
//	}
//	defer rows.Close()
//
//	var users []*pb.User
//	for rows.Next() {
//		var user pb.User
//		var role string
//		var createdAt, updatedAt time.Time
//
//		if err := rows.Scan(&user.Id, &user.Name, &user.Email, &role, &createdAt, &updatedAt); err != nil {
//			return nil, status.Errorf(codes.Internal, "failed to scan user: %v", err)
//		}
//
//		user.Role = pb.Role(pb.Role_value[role])
//		user.CreatedAt = timestamppb.New(createdAt)
//		user.UpdatedAt = timestamppb.New(updatedAt)
//		users = append(users, &user)
//	}
//
//	if err := rows.Err(); err != nil {
//		return nil, status.Errorf(codes.Internal, "error iterating through rows: %v", err)
//	}
//
//	logger.Info("users fetched")
//
//	return &pb.ListResponse{Users: users}, nil
//}

func main() {
	cfg := config.MustLoad()
	ctx := context.Background()

	pool, err := pgxpool.Connect(ctx, cfg.Database.PostgresDSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	userRepository := userRepo.NewRepository(pool)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPC.GRPCPort))
	if err != nil {
		logger.Fatal("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	pb.RegisterUserV1Server(s, &server{userRepository: userRepository})

	logger.Info("%v listening at %v", cfg.AppName, lis.Addr())

	if err = s.Serve(lis); err != nil {
		logger.Fatal("failed to serve: %v", err)
	}
}
