package mssql

import (
	"context"
	"github.com/jmoiron/sqlx"
	"gitlab/nefco/platform/codegen/generate/service/mssql/build"
	_query "gitlab/nefco/platform/codegen/generate/service/mssql/query"
)

func UpdateMany(ctx context.Context, result interface{}) (error) {
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

	sqlxQuery, args, err := sqlx.Named(query.Query(), query.Arg())
	if err != nil {
		return err
	}

	sqlxQuery, args, err = sqlx.In(sqlxQuery, args...)
	if err != nil {
		return err
	}

	sqlxQuery = tx.Rebind(sqlxQuery)

	_, err = tx.Exec(sqlxQuery, args...)
	if err != nil {
		return err
	}

	if err := Collection(ctx, result); err != nil {
		return err
	}

	return nil
}
