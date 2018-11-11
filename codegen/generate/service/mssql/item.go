package mssql

import (
	"context"
	"gitlab/nefco/platform/codegen/generate/service/mssql/query"
)

func Item(ctx context.Context, result interface{}) error {
	query := new(query.Select)

	if err := fillTable(ctx, query); err != nil {
		return err
	}

	if err := fillColumns(ctx, query); err != nil {
		return err
	}

	if err := fillConditions(ctx, query); err != nil {
		return err
	}

	logQuery(query)

	tx, err := Begin(ctx)
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareNamed(query.Query())
	if err != nil {
		return err
	}

	if err := stmt.Get(result, query.Arg()); err != nil {
		return err
	}

	return nil
}
