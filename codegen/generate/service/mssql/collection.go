package mssql

import (
	"context"
	"gitlab/nefco/platform/codegen/generate/service/mssql/build"
	"gitlab/nefco/platform/codegen/generate/service/mssql/query"

	"github.com/jmoiron/sqlx"
)

func Collection(ctx context.Context, result interface{}) error {
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

	if err := build.Pagination(ctx, query); err != nil {
		return err
	}

	if err := build.Order(ctx, query); err != nil {
		return err
	}

	if err := build.Trasher(ctx, query); err != nil {
		return err
	}

	if err := fillSoftDeleteFieldName(ctx, query); err != nil {
		return err
	}

	logQuery(query)

	tx, err := Begin(ctx)
	if err != nil {
		return err
	}

	_query, args, err := sqlx.Named(query.Query(), query.Arg())
	if err != nil {
		return err
	}

	_query, args, err = sqlx.In(_query, args...)
	if err != nil {
		return err
	}
	_query = tx.Rebind(_query)
	if err := tx.Select(result, _query, args...); err != nil {
		return err
	}

	return nil
}
