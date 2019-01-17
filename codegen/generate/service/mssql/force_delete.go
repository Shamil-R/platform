package mssql

import (
	"context"
	"gitlab/nefco/platform/codegen/generate/service/mssql/build"
	"gitlab/nefco/platform/codegen/generate/service/mssql/query"
)

func ForceDelete(ctx context.Context, result interface{}) error {
	if err := Item(ctx, result); err != nil {
		return err
	}

	query := query.NewForceDelete()

	if err := build.TableFromContext(ctx, query); err != nil {
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

	return nil
}
