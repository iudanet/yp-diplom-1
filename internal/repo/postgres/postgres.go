package postgres

import (
	"database/sql"
	"fmt"

	"github.com/iudanet/yp-diplom-1/internal/repo"
	"github.com/iudanet/yp-diplom-1/internal/repo/migrator"
	_ "github.com/lib/pq"
)

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(dsn string) (repo.Repositories, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err = migrator.Migrate(db); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	return &PostgresRepo{db: db}, nil
}
