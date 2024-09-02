package pg

import (
	"context"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/mikhailsoldatkin/auth/internal/repository/user/pg/filter"
	"github.com/mikhailsoldatkin/platform_common/pkg/db"
	"golang.org/x/crypto/bcrypt"

	"github.com/mikhailsoldatkin/auth/internal/customerrors"
	"github.com/mikhailsoldatkin/auth/internal/repository"
	"github.com/mikhailsoldatkin/auth/internal/repository/user/pg/converter"
	repoModel "github.com/mikhailsoldatkin/auth/internal/repository/user/pg/model"
	"github.com/mikhailsoldatkin/auth/internal/service/user/model"
)

const (
	tableUsers       = "users"
	tablePermissions = "permissions"
	columnID         = "id"
	columnUsername   = "username"
	columnEmail      = "email"
	columnRole       = "role"
	columnPassword   = "password"
	columnCreatedAt  = "created_at"
	columnUpdatedAt  = "updated_at"
	userEntity       = "user"
	columnEndpoint   = "endpoint"

	defaultPageSize = 10
)

var _ repository.UserRepository = (*repo)(nil)

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
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	builder := sq.Insert(tableUsers).
		PlaceholderFormat(sq.Dollar).
		Columns(
			columnUsername,
			columnEmail,
			columnRole,
			columnPassword,
			columnCreatedAt,
			columnUpdatedAt,
		).
		Values(user.Username, user.Email, user.Role, password, now, now).
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

// Get retrieves a user by given filter parameter from the database.
func (r *repo) Get(ctx context.Context, f filter.UserFilter) (*model.User, error) {
	err := f.Validate()
	if err != nil {
		return nil, err
	}

	builder := sq.Select(
		columnID,
		columnUsername,
		columnEmail,
		columnRole,
		columnPassword,
		columnCreatedAt,
		columnUpdatedAt,
	).
		From(tableUsers).
		PlaceholderFormat(sq.Dollar)

	var notFoundErr error

	if f.ID != nil {
		builder = builder.Where(sq.Eq{columnID: *f.ID})
		notFoundErr = customerrors.NewErrNotFound(userEntity, *f.ID)
	} else {
		builder = builder.Where(sq.Eq{columnUsername: *f.Username})
		notFoundErr = customerrors.NewErrNotFound(userEntity, *f.Username)
	}

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
			return nil, notFoundErr
		}
		return nil, err
	}

	return converter.FromRepoToService(&user), nil
}

// GetEndpointRoles retrieves roles associated with a specific endpoint from the database.
func (r *repo) GetEndpointRoles(ctx context.Context, endpoint string) ([]string, error) {
	builder := sq.Select(columnRole).
		From(tablePermissions).
		Where(sq.Eq{columnEndpoint: endpoint}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "user_repository.GetEndpointRoles",
		QueryRaw: query,
	}

	var roles []string
	err = r.db.DB().ScanAllContext(ctx, &roles, q, args...)
	if err != nil {
		return nil, err
	}

	return roles, nil
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

// CheckUsersExist checks if all users with the given IDs exist in the database.
// It returns an error if any of the provided IDs do not exist.
func (r *repo) CheckUsersExist(ctx context.Context, ids []int64) error {
	builder := sq.Select(columnID).
		From(tableUsers).
		Where(sq.Eq{columnID: ids}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "user_repository.CheckUsersExist",
		QueryRaw: query,
	}

	var dbIDs []int64
	err = r.db.DB().ScanAllContext(ctx, &dbIDs, q, args...)
	if err != nil {
		return err
	}

	existingIDMap := make(map[int64]bool, len(dbIDs))
	for _, id := range dbIDs {
		existingIDMap[id] = true
	}

	for _, id := range ids {
		if _, exists := existingIDMap[id]; !exists {
			return customerrors.NewErrNotFound(userEntity, id)
		}
	}

	return nil
}
