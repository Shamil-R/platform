package mssql

import (
	"context"
	"gitlab/nefco/platform/codegen/generate/service/mssql/build"
	_query "gitlab/nefco/platform/codegen/generate/service/mssql/query"
)

func Update(ctx context.Context, result interface{}) error {
	query := _query.NewUpdate()

	if err := build.TableFromSchema(ctx, query); err != nil {
		return err
	}

	data, err := build.ExtractArgument(ctx, "update")
	if err != nil {
		return err
	}

	if err := build.Value(data, query); err != nil {
		return err
	}

	if err := build.Updated(ctx, query); err != nil {
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

	if _, err := tx.NamedExec(query.Query(), query.Arg()); err != nil {
		return err
	}

	if err := Item(ctx, result); err != nil {
		return err
	}

	return nil
}
