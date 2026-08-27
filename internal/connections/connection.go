package connections

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func NewPool(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	if err := goose.SetDialect("postgres"); err != nil {
		pool.Close()
		return nil, err
	}

	db := stdlib.OpenDBFromPool(pool)

	if err := goose.Up(db, "migrations"); err != nil {
		db.Close()
		pool.Close()
		return nil, err
	}

	if err := db.Close(); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}
