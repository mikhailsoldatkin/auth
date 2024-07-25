package user

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/brianvoe/gofakeit"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mikhailsoldatkin/auth/internal/repository"
	"github.com/mikhailsoldatkin/auth/internal/repository/user/converter"
	"github.com/mikhailsoldatkin/auth/internal/repository/user/model"
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
