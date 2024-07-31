package log

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/mikhailsoldatkin/auth/internal/client/db"
	"github.com/mikhailsoldatkin/auth/internal/repository"
)

const (
	tableUsersLogs = "users_logs"
	columnDetails  = "details"
	columnUserID   = "user_id"
)

type logRepo struct {
	db db.Client
}

// NewRepository creates a new instance of the log repository.
func NewRepository(db db.Client) repository.LogRepository {
	return &logRepo{db: db}
}

// Log logs a database operation.
func (r *logRepo) Log(ctx context.Context, userID int64, details string) error {
	// если апдейтим без указания id, будет ошибка, возможно эта проверка должна быть не здесь
	if userID == 0 {
		return nil
	}

	builder := sq.Insert(tableUsersLogs).
		PlaceholderFormat(sq.Dollar).
		Columns(columnUserID, columnDetails).
		Values(userID, details)

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build SQL query: %w", err)
	}

	q := db.Query{
		Name:     "log_repository.Log",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("failed to execute SQL query: %w", err)
	}

	return nil
}
