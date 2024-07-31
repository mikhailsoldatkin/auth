package user

import (
	"context"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/mikhailsoldatkin/auth/internal/client/db"
	"github.com/mikhailsoldatkin/auth/internal/customerrors"
	"github.com/mikhailsoldatkin/auth/internal/repository"
	"github.com/mikhailsoldatkin/auth/internal/repository/user/converter"
	modelRepo "github.com/mikhailsoldatkin/auth/internal/repository/user/model"
	"github.com/mikhailsoldatkin/auth/internal/service/user/model"
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
	db db.Client
}

// NewRepository creates a new instance of the user repository.
func NewRepository(db db.Client) repository.UserRepository {
	return &repo{db: db}
}

// Create inserts a new user into the database.
func (r *repo) Create(ctx context.Context, user *model.User) (int64, error) {
	builder := sq.Insert(tableUsers).
		PlaceholderFormat(sq.Dollar).
		Columns(
			columnName,
			columnEmail,
			columnRole,
		).
		Values(user.Name, user.Email, user.Role).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "user_repository.Create",
		QueryRaw: query,
	}

	var id int64
	err = r.db.DB().ScanOneContext(ctx, &id, q, args...)
	if err != nil {
		return 0, status.Errorf(codes.Internal, "failed to execute query: %v", err)
	}

	return id, nil
}

// Get retrieves a user by ID from the database.
func (r *repo) Get(ctx context.Context, id int64) (*model.User, error) {
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

	q := db.Query{
		Name:     "user_repository.Get",
		QueryRaw: query,
	}

	var user modelRepo.User
	err = r.db.DB().ScanOneContext(ctx, &user, q, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, customerrors.ErrNotFound
		}
		return nil, err
	}

	return converter.ToServiceFromRepo(&user), nil
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

	q := db.Query{
		Name:     "user_repository.Delete",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return nil
}

// Update modifies an existing user in the database.
func (r *repo) Update(ctx context.Context, req *pb.UpdateRequest) error {
	updateFields := make(map[string]any)
	updateFields[columnUpdatedAt] = time.Now()

	if req.GetName() != nil {
		updateFields[columnName] = req.GetName().GetValue()
	}
	if req.GetEmail() != nil {
		updateFields[columnEmail] = req.GetEmail().GetValue()
	}
	if req.GetRole().String() != "" {
		updateFields[columnRole] = req.GetRole().String()
	}

	builder := sq.Update(tableUsers).
		SetMap(updateFields).
		Where(sq.Eq{columnID: req.GetId()}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "user_repository.Update",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return customerrors.ErrNotFound
		}
		return err
	}

	return nil
}

// List retrieves a list of users from the database.
func (r *repo) List(ctx context.Context, req *pb.ListRequest) ([]*model.User, error) {
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

	q := db.Query{
		Name:     "user_repository.List",
		QueryRaw: query,
	}

	var usersRepo []*modelRepo.User
	err = r.db.DB().ScanAllContext(ctx, &usersRepo, q, args...)
	if err != nil {
		return nil, err
	}

	return converter.ToServiceFromRepoList(usersRepo), nil
}
