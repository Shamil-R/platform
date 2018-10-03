package mssql

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/99designs/gqlgen/graphql"
)

type SQLContext struct {
	db *sqlx.DB
	tx *sqlx.Tx
}

func NewContext(db *sqlx.DB) *SQLContext {
	return &SQLContext{db: db}
}

func (c *SQLContext) IsTx() bool {
	return c.tx != nil
}

func (c *SQLContext) Begin() (*sqlx.Tx, error) {
	if c.IsTx() {
		return c.tx, nil
	}
	tx, err := c.db.Beginx()
	if err != nil {
		return nil, err
	}
	c.tx = tx
	return tx, nil
}

func (c *SQLContext) Commit() error {
	if !c.IsTx() {
		return nil
	}
	if err := c.tx.Commit(); err != nil {
		return err
	}
	c.tx = nil
	return nil
}

func (c *SQLContext) Rollback() error {
	if !c.IsTx() {
		return nil
	}
	if err := c.tx.Rollback(); err != nil {
		return err
	}
	c.tx = nil
	return nil
}

type key int

const ctxKey key = 1

func withContext(ctx context.Context, sqlCtx *SQLContext) context.Context {
	return context.WithValue(ctx, ctxKey, sqlCtx)
}

func getContext(ctx context.Context) *SQLContext {
	val := ctx.Value(ctxKey)
	if val == nil {
		return nil
	}
	return val.(*SQLContext)
}

func Begin(ctx context.Context) (*sqlx.Tx, error) {
	return getContext(ctx).Begin()
}

func middleware(db *sqlx.DB) graphql.RequestMiddleware {
	return func(ctx context.Context, next func(ctx context.Context) []byte) []byte {
		sqlCtx := NewContext(db)

		res := next(withContext(ctx, sqlCtx))

		if sqlCtx.IsTx() {
			reqCtx := graphql.GetRequestContext(ctx)
			if len(reqCtx.Errors) == 0 {
				if err := sqlCtx.Commit(); err != nil {
					panic(err)
				}
			} else {
				if err := sqlCtx.Rollback(); err != nil {
					panic(err)
				}
			}
		}

		return res
	}
}
