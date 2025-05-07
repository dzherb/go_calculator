package storage

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Retriever interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type TransactionIssuer interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

type Connection interface {
	Retriever
	TransactionIssuer
}
