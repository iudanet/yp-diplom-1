package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/iudanet/yp-diplom-1/internal/pkg/retry"
	"github.com/iudanet/yp-diplom-1/internal/repo"
	"github.com/iudanet/yp-diplom-1/internal/repo/migrator"
	_ "github.com/lib/pq"
)

var _ repo.Repositories = (*postgresRepo)(nil)

const (
	timeoutPing    = 5 * time.Second
	timeoutMigrate = 20 * time.Second
)

type postgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(ctx context.Context, dsn string) (repo.Repositories, error) {
	var db *sql.DB
	var err error

	// Используем retry для подключения к базе данных
	err = retry.WithRetry(func() error {
		db, err = sql.Open("postgres", dsn)
		if err != nil {
			return err
		}

		ctxPing, cancel := context.WithTimeout(ctx, timeoutPing)
		defer cancel()

		if err := db.PingContext(ctxPing); err != nil {
			_ = db.Close() // закрываем соединение при ошибке
			return err
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &postgresRepo{db: db}, nil
}

func (r *postgresRepo) Migrate(ctx context.Context) error {
	err := retry.WithRetry(func() error {
		ctxMigrate, cancel := context.WithTimeout(ctx, timeoutMigrate)
		defer cancel()
		return migrator.Migrate(ctxMigrate, r.db)
	})
	if err != nil {
		_ = r.db.Close()
		return fmt.Errorf("migration failed: %w", err)
	}
	return nil
}
