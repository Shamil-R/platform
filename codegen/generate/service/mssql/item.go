package mssql

import (
	"context"
	"gitlab/nefco/platform/codegen/generate/service/mssql/build"
	"gitlab/nefco/platform/codegen/generate/service/mssql/query"
)

func Item(ctx context.Context, result interface{}) error {
	query := query.NewSelect()

	if err := build.TableFromSelection(ctx, query); err != nil {
		return err
	}

	if err := build.ColumnsFromSelection(ctx, query); err != nil {
		return err
	}

	if err := build.Conditions(ctx, query); err != nil {
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
