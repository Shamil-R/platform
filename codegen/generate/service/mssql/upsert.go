package mssql

import (
	"context"
	"database/sql"
)

func Upsert(ctx context.Context, result interface{}) error {
	err := Item(ctx, result)
	if err != sql.ErrNoRows && err != nil {
		return err
	}

	if err != sql.ErrNoRows {
		return Update(ctx, result)
	} else {
		return Create(ctx, result)
	}


	return nil
}
