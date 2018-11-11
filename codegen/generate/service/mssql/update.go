package mssql

import (
	"context"
	"gitlab/nefco/platform/codegen/generate/service/mssql/query"
)

func Update(ctx context.Context, result interface{}) error {
	query := new(query.Update)

	if err := fillTable(ctx, query); err != nil {
		return err
	}

	if err := fillValues(ctx, query); err != nil {
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

	if _, err := tx.NamedExec(query.Query(), query.Arg()); err != nil {
		return err
	}

	if err := Item(ctx, result); err != nil {
		return err
	}

	return nil
}
