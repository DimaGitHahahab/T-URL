package postgres

import (
	"context"
	"errors"
	"fmt"

	"storage/internal/repository"
	"storage/internal/repository/postgres/queries"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresRepo struct {
	*queries.Queries // SQL queries
	pool             *pgxpool.Pool
}

func New(pgxPool *pgxpool.Pool) repository.Repository {
	r := &postgresRepo{
		Queries: queries.New(pgxPool),
		pool:    pgxPool,
	}

	return r
}

func SetupPgxPool(ctx context.Context, DbURL string) (*pgxpool.Pool, error) {
	pgxConfig, err := pgxpool.ParseConfig(DbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pgx config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create new pgx pool with config: %w", err)
	}

	return pool, nil
}

func ProcessMigration(migrationURL string, dbSource string) error {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		return fmt.Errorf("failed to create new migration: %w", err)
	}

	if err = migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to migrate up: %w", err)
	}
	defer migration.Close()

	return nil
}
