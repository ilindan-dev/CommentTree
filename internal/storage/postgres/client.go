package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/ilindan-dev/CommentTree/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrParsePoolConfig indicates an error parsing the database pool configuration.
var ErrParsePoolConfig = errors.New("error parsing pool config")

// ErrCreatePool indicates an error creating the database pool.
var ErrCreatePool = errors.New("error creating database pool")

// ErrPingDB indicates an error pinging the database.
var ErrPingDB = errors.New("error pinging database")

// NewClient creates a new PostgreSQL client with the given configuration.
func NewClient(ctx context.Context, cfg config.DBConfig) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrParsePoolConfig, err)
	}

	if cfg.Pool.MaxOpenConns > 0 {
		poolCfg.MaxConns = cfg.Pool.MaxOpenConns
	}
	if cfg.Pool.MaxIdleConns > 0 {
		poolCfg.MinConns = cfg.Pool.MaxIdleConns
	}
	if cfg.Pool.MaxIdleTime > 0 {
		poolCfg.MaxConnIdleTime = cfg.Pool.MaxIdleTime
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCreatePool, err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrPingDB, err)
	}

	return pool, nil
}
