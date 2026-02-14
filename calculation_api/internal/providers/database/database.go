package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	connPool *pgxpool.Pool
	*Queries
}

func NewDatabase(connPool *pgxpool.Pool) *Database {
	return &Database{
		connPool: connPool,
		Queries:  New(connPool),
	}
}

func (d *Database) Close() {
	d.connPool.Close()
}

func (d *Database) RawQuery(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return d.connPool.Query(ctx, sql, args...)
}

func (d *Database) RawQueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return d.connPool.QueryRow(ctx, sql, args...)
}

func (d *Database) Ping(ctx context.Context) error {
	return d.connPool.Ping(ctx)
}
