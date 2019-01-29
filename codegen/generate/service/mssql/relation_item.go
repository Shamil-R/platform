package mssql

import (
	"context"
	"github.com/jmoiron/sqlx"

	"gitlab/nefco/platform/codegen/generate/service/mssql/build"
	"gitlab/nefco/platform/codegen/generate/service/mssql/query"
)

//todo deleted_at is null
func RelationItem(ctx context.Context, objID int, result interface{}) error {
	child := query.NewSelect()

	if err := build.TableFromSelection(ctx, child); err != nil {
		return err
	}

	if err := build.ColumnsFromSelection(ctx, child); err != nil {
		return err
	}

	if err := build.Conditions(ctx, child); err != nil {
		return err
	}

	if err := build.Pagination(ctx, child); err != nil {
		return err
	}

	col, err := getRelationColumn(ctx)
	if err != nil {
		return err
	}

	child.AddСondition(col, "eq", objID)

	logQuery(child)

	tx, err := Begin(ctx)
	if err != nil {
		return err
	}

	_query, args, err := sqlx.Named(child.Query(), child.Arg())
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
