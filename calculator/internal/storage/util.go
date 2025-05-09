package storage

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

func WithTransaction(
	ctx context.Context,
	db Connection,
	f func(pgx.Tx) error,
) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err = tx.Rollback(ctx); err != nil &&
			!errors.Is(err, pgx.ErrTxClosed) {
			slog.Error(err.Error())
		}
	}()

	if err = f(tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
