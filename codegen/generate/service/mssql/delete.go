package mssql

import (
	"context"
	"gitlab/nefco/platform/codegen/generate/service/mssql/build"
	"gitlab/nefco/platform/codegen/generate/service/mssql/query"
	"time"
)

func Delete(ctx context.Context, result interface{}) error {
	query := query.NewUpdate()

	if err := build.Table(ctx, query); err != nil {
		return err
	}

	dirName := "softDelete"
	argName := "deleteField"
	fieldName, err := getDefaultValues(ctx, dirName, argName)
	if err != nil {
		return err
	}
	query.AddValue(fieldName, time.Now())

	if err := build.Conditions(ctx, query); err != nil {
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
