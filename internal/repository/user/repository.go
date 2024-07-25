package user

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/brianvoe/gofakeit"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mikhailsoldatkin/auth/internal/repository"
	"github.com/mikhailsoldatkin/auth/internal/repository/user/converter"
	"github.com/mikhailsoldatkin/auth/internal/repository/user/model"
	"github.com/mikhailsoldatkin/auth/internal/utils"
	pb "github.com/mikhailsoldatkin/auth/pkg/user_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	tableUsers      = "users"
	columnID        = "id"
	columnName      = "name"
	columnEmail     = "email"
	columnRole      = "role"
	columnCreatedAt = "created_at"
	columnUpdatedAt = "updated_at"

	defaultPageSize = 10
)

type repo struct {
	db *pgxpool.Pool
}

// NewRepository creates a new instance of the user repository.
func NewRepository(db *pgxpool.Pool) repository.UserRepository {
	return &repo{db: db}
}

// checkUserExists checks if user with given ID exists in database and returns an error if it doesn't.
func (r *repo) checkUserExists(ctx context.Context, userID int64) error {
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE %s=$1)", tableUsers, columnID)
	err := r.db.QueryRow(ctx, query, userID).Scan(&exists)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to check user existence: %v", err)
	}
	if !exists {
		return status.Errorf(codes.NotFound, "user with ID %d not found", userID)
	}
	return nil
}

// Create inserts a new user into the database.
func (r *repo) Create(ctx context.Context, user *pb.User) (int64, error) {
	builder := sq.Insert(tableUsers).
		PlaceholderFormat(sq.Dollar).
		Columns(
			columnName,
			columnEmail,
			columnRole,
		).
		Values(gofakeit.Name(), gofakeit.Email(), user.Role.String()).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	var id int
	err = r.db.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return int64(id), nil
}

// Get retrieves a user by ID from the database.
func (r *repo) Get(ctx context.Context, id int64) (*pb.User, error) {
	if err := r.checkUserExists(ctx, id); err != nil {
		return nil, err
	}
	builder := sq.Select(
		columnID,
		columnName,
		columnEmail,
		columnRole,
		columnCreatedAt,
		columnUpdatedAt,
	).
		From(tableUsers).
		Where(sq.Eq{columnID: id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	var user model.User
	err = r.db.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return converter.ToUserFromRepo(&user), nil
}

// Delete removes a user by ID from the database.
func (r *repo) Delete(ctx context.Context, id int64) error {
	builder := sq.Delete(tableUsers).
		Where(sq.Eq{columnID: id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

// Update modifies an existing user in the database.
func (r *repo) Update(ctx context.Context, req *pb.UpdateRequest) error {
	if err := r.checkUserExists(ctx, req.GetId()); err != nil {
		return err
	}

	updateFields := make(map[string]any)

	if req.GetName() != nil {
		updateFields[columnName] = req.GetName().GetValue()
	}
	if req.GetEmail() != nil {
		email := req.GetEmail().GetValue()
		if !utils.ValidateEmail(email) {
			return status.Errorf(codes.InvalidArgument, "invalid email format: %v", email)
		}
		updateFields[columnEmail] = email
	}
	if req.GetRole().String() != "" {
		updateFields[columnRole] = req.GetRole().String()
	}

	if len(updateFields) == 0 {
		return status.Errorf(codes.InvalidArgument, "no fields to update")
	}

	updateFields[columnUpdatedAt] = time.Now()

	builder := sq.Update(tableUsers).
		SetMap(updateFields).
		Where(sq.Eq{columnID: req.GetId()}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

// List retrieves a list of users from the database.
func (r *repo) List(ctx context.Context, req *pb.ListRequest) ([]*pb.User, error) {
	limit := int(req.GetLimit())
	if limit <= 0 {
		limit = defaultPageSize
	}

	builder := sq.Select("*").
		From(tableUsers).
		OrderBy(columnID).
		Limit(uint64(limit)).
		Offset(uint64(req.GetOffset())).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	usersPb := make([]*pb.User, 0, len(users))
	for _, user := range users {
		usersPb = append(usersPb, converter.ToUserFromRepo(user))
	}

	return usersPb, nil
}
