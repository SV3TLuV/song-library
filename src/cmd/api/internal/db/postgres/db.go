package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

func NewDB(ctx context.Context, conn string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(conn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse config")
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "init pool")
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, errors.Wrap(err, "ping pool")
	}

	return pool, nil
}
