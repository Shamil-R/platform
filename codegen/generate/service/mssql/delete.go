package mssql

import (
	"context"
	"gitlab/nefco/platform/codegen/generate/service/mssql/build"
	"gitlab/nefco/platform/codegen/generate/service/mssql/query"
)

func Delete(ctx context.Context, result interface{}) error {
	query := query.NewUpdate()

	data, err := build.ExtractArgument(ctx, "where")
	if err != nil {
		return err
	}

	if err := build.TableFromValue(data, query); err != nil {
		return err
	}

	if err := build.SoftDelete(data, query); err != nil {
		return err
	}

	if err := build.ConditionsFromValue(data, query); err != nil {
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
