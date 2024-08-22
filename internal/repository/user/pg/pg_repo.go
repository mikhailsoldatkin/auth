package pg

import (
	"context"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/mikhailsoldatkin/platform_common/pkg/db"

	"github.com/mikhailsoldatkin/auth/internal/customerrors"
	"github.com/mikhailsoldatkin/auth/internal/repository"
	"github.com/mikhailsoldatkin/auth/internal/repository/user/pg/converter"
	repoModel "github.com/mikhailsoldatkin/auth/internal/repository/user/pg/model"
	"github.com/mikhailsoldatkin/auth/internal/service/user/model"
)

const (
	tableUsers      = "users"
	columnID        = "id"
	columnName      = "name"
	columnEmail     = "email"
	columnRole      = "role"
	columnCreatedAt = "created_at"
	columnUpdatedAt = "updated_at"
	userEntity      = "user"

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
	now := time.Now()
	builder := sq.Insert(tableUsers).
		PlaceholderFormat(sq.Dollar).
		Columns(
			columnName,
			columnEmail,
			columnRole,
			columnCreatedAt,
			columnUpdatedAt,
		).
		Values(user.Name, user.Email, user.Role, now, now).
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
		return 0, err
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

	var user repoModel.User
	err = r.db.DB().ScanOneContext(ctx, &user, q, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, customerrors.NewErrNotFound(userEntity, id)
		}
		return nil, err
	}

	return converter.FromRepoToService(&user), nil
}

// Delete removes a user from the database by ID.
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
func (r *repo) Update(ctx context.Context, updates *model.User) error {
	updateFields := converter.FromServiceToRepoUpdate(updates)

	builder := sq.Update(tableUsers).
		SetMap(updateFields).
		Where(sq.Eq{columnID: updates.ID}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "user_repository.Update",
		QueryRaw: query,
	}

	result, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return customerrors.NewErrNotFound(userEntity, updates.ID)
	}

	return nil
}

// List retrieves a list of users from the database.
func (r *repo) List(ctx context.Context, limit, offset int64) ([]*model.User, error) {
	if limit <= 0 {
		limit = defaultPageSize
	}
	if offset < 0 {
		offset = 0
	}

	builder := sq.Select("*").
		From(tableUsers).
		OrderBy(columnID).
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "user_repository.List",
		QueryRaw: query,
	}

	var repoUsers []*repoModel.User
	err = r.db.DB().ScanAllContext(ctx, &repoUsers, q, args...)
	if err != nil {
		return nil, err
	}

	return converter.FromRepoToServiceList(repoUsers), nil
}
