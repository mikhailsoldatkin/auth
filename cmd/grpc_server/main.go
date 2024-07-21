package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"regexp"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/brianvoe/gofakeit"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mikhailsoldatkin/auth/internal/config"
	"github.com/mikhailsoldatkin/auth/internal/logger"
	pb "github.com/mikhailsoldatkin/auth/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type server struct {
	pb.UnimplementedUserV1Server
	pool *pgxpool.Pool
}

type pageToken struct {
	LastID int64 `json:"last_id"`
}

func encodePageToken(token pageToken) (string, error) {
	tokenBytes, err := json.Marshal(token)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(tokenBytes), nil
}

func decodePageToken(tokenStr string) (pageToken, error) {
	tokenBytes, err := base64.URLEncoding.DecodeString(tokenStr)
	if err != nil {
		return pageToken{}, err
	}
	var token pageToken
	err = json.Unmarshal(tokenBytes, &token)
	if err != nil {
		return pageToken{}, err
	}
	return token, nil
}

// ensureUserExists checks if a user with given ID exists in database and returns an error if it doesn't.
func (s *server) ensureUserExists(ctx context.Context, userID int64) error {
	var exists bool
	err := s.pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE id=$1)", userID).Scan(&exists)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to check user existence: %v", err)
	}
	if !exists {
		return status.Errorf(codes.NotFound, "user with ID %d not found", userID)
	}
	return nil
}

// roleToString converts a pb.Role to its corresponding string representation.
func roleToString(role pb.Role) (string, error) {
	roleMap := map[pb.Role]string{
		pb.Role_UNKNOWN: "UNKNOWN",
		pb.Role_USER:    "USER",
		pb.Role_ADMIN:   "ADMIN",
	}

	roleStr, ok := roleMap[role]
	if !ok {
		return "", fmt.Errorf("unknown role: %v", role)
	}
	return roleStr, nil
}

// stringToRole converts a string representation of a role to pb.Role.
func stringToRole(roleStr string) (pb.Role, error) {
	roleMap := map[string]pb.Role{
		"UNKNOWN": pb.Role_UNKNOWN,
		"USER":    pb.Role_USER,
		"ADMIN":   pb.Role_ADMIN,
	}

	role, ok := roleMap[roleStr]
	if !ok {
		return pb.Role_UNKNOWN, fmt.Errorf("unknown role string: %s", roleStr)
	}
	return role, nil
}

// validateEmail checks if the given email address is in a valid format.
func validateEmail(email string) bool {
	const emailRegex = `^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

// validatePassword provides simple password validation.
func validatePassword(password, passwordConfirm string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	if password != passwordConfirm {
		return errors.New("passwords do not match")
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

	builder := sq.Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns(
			"name",
			"email",
			"role",
		).
		//Values(req.GetName(), req.GetEmail(), req.GetRole().String()).
		Values(gofakeit.Name(), gofakeit.Email(), req.GetRole().String()).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to build query: %v", err)
	}

	var userID int
	err = s.pool.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	logger.Info("user with ID %d created", userID)

	return &pb.CreateResponse{Id: int64(userID)}, nil
}

// Get retrieves user data by ID.
func (s *server) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	builder := sq.Select(
		"id",
		"name",
		"email",
		"role",
		"created_at",
		"updated_at",
	).
		From("users").
		Where(sq.Eq{"id": req.GetId()}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to build query: %v", err)
	}

	var user pb.User
	var roleStr string
	var createdAt, updatedAt time.Time

	err = s.pool.QueryRow(ctx, query, args...).Scan(&user.Id, &user.Name, &user.Email, &roleStr, &createdAt, &updatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "user with ID %d not found", req.GetId())
		}
		return nil, status.Errorf(codes.Internal, "failed to retrieve user: %v", err)
	}

	role, err := stringToRole(roleStr)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert role: %v", err)
	}

	user.Role = role
	user.CreatedAt = timestamppb.New(createdAt)
	user.UpdatedAt = timestamppb.New(updatedAt)

	logger.Info("user data retrieved %v", &user)

	return &pb.GetResponse{User: &user}, nil
}

// Update modifies user data.
func (s *server) Update(ctx context.Context, req *pb.UpdateRequest) (*emptypb.Empty, error) {
	if err := s.ensureUserExists(ctx, req.GetId()); err != nil {
		return nil, err
	}

	updateFields := make(map[string]any)

	if req.GetName() != nil {
		updateFields["name"] = req.GetName().GetValue()
	}
	if req.GetEmail() != nil {
		email := req.GetEmail().GetValue()
		if !validateEmail(email) {
			return nil, status.Errorf(codes.InvalidArgument, "invalid email format: %v", email)
		}
		updateFields["email"] = email
	}
	if req.GetRole() != pb.Role_UNKNOWN {
		roleStr, err := roleToString(req.GetRole())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid role value: %v", req.GetRole())
		}
		updateFields["role"] = roleStr
	}

	if len(updateFields) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "no fields to update")
	}

	updateFields["updated_at"] = time.Now()

	builder := sq.Update("users").
		SetMap(updateFields).
		Where(sq.Eq{"id": req.GetId()}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to build update query: %v", err)
	}

	_, err = s.pool.Exec(ctx, query, args...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}

	logger.Info("user with ID &d updated", req.GetId())

	return &emptypb.Empty{}, nil
}

// Delete removes a user by ID.
func (s *server) Delete(ctx context.Context, req *pb.DeleteRequest) (*emptypb.Empty, error) {
	if err := s.ensureUserExists(ctx, req.GetId()); err != nil {
		return nil, err
	}

	builder := sq.Delete("users").
		Where(sq.Eq{"id": req.GetId()}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to build delete query: %v", err)
	}

	_, err = s.pool.Exec(ctx, query, args...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}

	logger.Info("user with ID %d deleted", int(req.GetId()))

	return &emptypb.Empty{}, nil
}

// List lists users with pagination support.
func (s *server) List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	const defaultPageSize = 10
	pageSize := int(req.GetPageSize())
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}

	var lastID int64
	if req.GetPageToken() != "" {
		token, err := decodePageToken(req.GetPageToken())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid page token: %v", err)
		}
		lastID = token.LastID
	}

	builder := sq.Select("*").
		From("users").
		Where(sq.Gt{"id": lastID}).
		OrderBy("id").
		Limit(uint64(pageSize + 1)).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to build query: %v", err)
	}

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list users: %v", err)
	}
	defer rows.Close()

	var users []*pb.User
	for rows.Next() {
		var user pb.User
		var roleStr string
		var createdAt, updatedAt time.Time

		if err := rows.Scan(&user.Id, &user.Name, &user.Email, &roleStr, &createdAt, &updatedAt); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to scan user: %v", err)
		}

		role, err := stringToRole(roleStr)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to convert role: %v", err)
		}

		user.Role = role
		user.CreatedAt = timestamppb.New(createdAt)
		user.UpdatedAt = timestamppb.New(updatedAt)
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "error iterating through rows: %v", err)
	}

	var nextPageToken string
	if len(users) > pageSize {
		nextPageToken, err = encodePageToken(pageToken{LastID: users[pageSize].Id})
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to create page token: %v", err)
		}
		users = users[:pageSize]
	}

	logger.Info("users fetched")

	return &pb.ListResponse{
		Users:         users,
		NextPageToken: nextPageToken,
	}, nil
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

	logger.Info("%v server listening at %v", cfg.AppName, lis.Addr())

	if err = s.Serve(lis); err != nil {
		logger.Fatal("failed to serve: %v", err)
	}
}
