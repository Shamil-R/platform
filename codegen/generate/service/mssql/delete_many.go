package mssql

import (
	"context"
	"github.com/jmoiron/sqlx"
	"gitlab/nefco/platform/codegen/generate/service/mssql/build"
	"gitlab/nefco/platform/codegen/generate/service/mssql/query"
)

func DeleteMany(ctx context.Context, result interface{}) (error) {
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

	if err := build.Conditions(ctx, query); err != nil {
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

	_, err = tx.Exec(_query, args...)
	if err != nil {
		return err
	}

	//result, err := rows.RowsAffected()
	//if err != nil {
	//	return err
	//}

	return nil
}
