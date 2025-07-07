package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/iudanet/yp-diplom-1/internal/repo"
	"github.com/iudanet/yp-diplom-1/internal/repo/migrator"
	"github.com/iudanet/yp-diplom-1/internal/retry"
	_ "github.com/lib/pq"
)

type PostgresRepo struct {
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

		ctxPing, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		if err := db.PingContext(ctxPing); err != nil {
			db.Close() // закрываем соединение при ошибке
			return err
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Применяем миграции
	err = retry.WithRetry(func() error {
		return migrator.Migrate(db)
	})
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	return &PostgresRepo{db: db}, nil
}
