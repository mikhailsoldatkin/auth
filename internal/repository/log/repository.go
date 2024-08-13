package log

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/mikhailsoldatkin/auth/internal/customerrors"
	"github.com/mikhailsoldatkin/auth/internal/repository"
	"github.com/mikhailsoldatkin/platform_common/pkg/db"
)

const (
	tableUsersLogs = "users_logs"
	columnDetails  = "details"
	columnUserID   = "user_id"
	userEntity     = "user"
)

type logRepo struct {
	db db.Client
}

// NewRepository creates a new instance of the log repository.
func NewRepository(db db.Client) repository.LogRepository {
	return &logRepo{db: db}
}

// Log logs a database operation on entity.
func (r *logRepo) Log(ctx context.Context, id int64, details string) error {
	builder := sq.Insert(tableUsersLogs).
		PlaceholderFormat(sq.Dollar).
		Columns(columnUserID, columnDetails).
		Values(id, details)

	if id == 0 {
		builder = sq.Insert(tableUsersLogs).
			PlaceholderFormat(sq.Dollar).
			Columns(columnDetails).
			Values(details)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "log_repository.Log",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if errors.Is(err, pgx.ErrNoRows) {
		return customerrors.NewErrNotFound(userEntity, id)
	}

	return nil
}
