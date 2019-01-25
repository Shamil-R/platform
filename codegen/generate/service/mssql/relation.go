package mssql

import (
	"context"
	"gitlab/nefco/platform/codegen/generate/service/mssql/build"
	"gitlab/nefco/platform/codegen/generate/service/mssql/query"

	"github.com/jmoiron/sqlx"
)

func Relation(ctx context.Context, objID int, result interface{}) error {
	//todo make select parent object
	//parent := query.NewSelect()




	q := query.NewSelect()

	if err := build.TableFromSelection(ctx, q); err != nil {
		return err
	}

	if err := build.ColumnsFromSelection(ctx, q); err != nil {
		return err
	}

	if err := build.Conditions(ctx, q); err != nil {
		return err
	}

	if err := build.Pagination(ctx, q); err != nil {
		return err
	}

	col, err := getRelationColumn(ctx)
	if err != nil {
		return err
	}

	q.AddСondition(col, "eq", objID)

	logQuery(q)

	tx, err := Begin(ctx)
	if err != nil {
		return err
	}

	_query, args, err := sqlx.Named(q.Query(), q.Arg())
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
