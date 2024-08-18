package database

import (
	"context"
	"database/sql"
	"fmt"
)

type txCtxKey struct{}

func TxFromContext(ctx context.Context) *sql.Tx {
	tx, _ := ctx.Value(txCtxKey{}).(*sql.Tx)
	return tx
}

func (db *Transactor) LoadClient(ctx context.Context) DBCaller {
	tx := TxFromContext(ctx)
	if tx != nil {
		return tx
	}
	return db.db
}

func NewTxContext(parent context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(parent, txCtxKey{}, tx)
}

type Transactor struct {
	db *sql.DB
}

func NewTransactor(DB *sql.DB) *Transactor {
	return &Transactor{db: DB}
}

type DBCaller interface {
	Exec(query string, args ...any) (sql.Result, error)
	QueryRow(query string, args ...any) *sql.Row
}

func (db *Transactor) RunInTx(ctx context.Context, f func(context.Context) error) (err error) {
	tx := TxFromContext(ctx)
	if tx != nil {
		return f(ctx)
	}

	tx, err = db.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}

	defer func() {
		if err := tx.Rollback(); err != nil {
			fmt.Printf("rollback error: %v", err)
		}
	}()

	if err = f(NewTxContext(ctx, tx)); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}
