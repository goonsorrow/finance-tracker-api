package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type TransactorPostgres struct {
	db *sqlx.DB
}

type TxKey struct{}

func NewTransactorPostgres(db *sqlx.DB) *TransactorPostgres {
	return &TransactorPostgres{db: db}
}

func (t *TransactorPostgres) GetExecutor(ctx context.Context) sqlx.ExtContext {
	if tx := GetTxFromContext(ctx); tx != nil {
		return tx
	}
	return t.db
}

func (t *TransactorPostgres) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {

	tx, err := t.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("[TransactorPostgres.WithinTransaction] begin tx: %w", err)
	}
	defer tx.Rollback()

	ctxKey := context.WithValue(ctx, TxKey{}, tx)

	if err := fn(ctxKey); err != nil {
		return fmt.Errorf("[TransactorPostgres.WithinTransaction] tx failed: %w", err)
	}

	return tx.Commit()
}

func GetTxFromContext(ctx context.Context) *sqlx.Tx {
	if tx, ok := ctx.Value(TxKey{}).(*sqlx.Tx); ok {
		return tx
	}
	return nil
}
