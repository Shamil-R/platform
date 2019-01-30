package mssql

import (
	"context"
	"github.com/jmoiron/sqlx"
	"gitlab/nefco/platform/codegen/generate/service/mssql/build"
	_query "gitlab/nefco/platform/codegen/generate/service/mssql/query"
)

//todo deleted_at is null
func RelationItem(ctx context.Context, objID int, result interface{}) error {
	child := _query.NewSelect()

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

	if err := build.RelationCondition(ctx, child, objID); err != nil {
		return err
	}

	logQuery(child)

	tx, err := Begin(ctx)
	if err != nil {
		return err
	}

	sqlxQuery, args, err := sqlx.Named(child.Query(), child.Arg())
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
