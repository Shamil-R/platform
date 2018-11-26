package mssql

import (
	"context"
	"github.com/jmoiron/sqlx"
	"gitlab/nefco/platform/codegen/generate/service/mssql/build"
	"gitlab/nefco/platform/codegen/generate/service/mssql/query"
)

func Relation(ctx context.Context, objID int, result interface{}) error {
	query := query.NewSelect()

	if err := fillTable(ctx, query); err != nil {
		return err
	}

	if err := fillColumns(ctx, query); err != nil {
		return err
	}

	if err := build.Conditions(ctx, query); err != nil {
		return err
	}

	if err := build.Pagination(ctx, query); err != nil {
		return err
	}

	col, err := getRelationColumn(ctx)
	if err != nil {
		return err
	}

	query.AddСondition(col, "eq", objID)

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
