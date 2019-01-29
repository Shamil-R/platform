package mssql

import (
	"context"
	"github.com/jmoiron/sqlx"
	"gitlab/nefco/platform/codegen/generate/service/mssql/build"
	_query "gitlab/nefco/platform/codegen/generate/service/mssql/query"
)

func Collection(ctx context.Context, result interface{}) error {
	query := _query.NewSelect()

	if err := build.TableFromSchema(ctx, query); err != nil {
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

	logQuery(query)

	tx, err := Begin(ctx)
	if err != nil {
		return err
	}

	sqlxQuery, args, err := sqlx.Named(query.Query(), query.Arg())
	if err != nil {
		return err
	}

	sqlxQuery, args, err = sqlx.In(sqlxQuery, args...)
	if err != nil {
		return err
	}
	sqlxQuery = tx.Rebind(sqlxQuery)
	if err := tx.Select(result, sqlxQuery, args...); err != nil {
		return err
	}

	return nil
}
