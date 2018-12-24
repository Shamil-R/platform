package mssql

import (
	"context"
	"github.com/jmoiron/sqlx"
	"gitlab/nefco/platform/codegen/generate/service/mssql/build"
	"gitlab/nefco/platform/codegen/generate/service/mssql/query"
)

func ForceDeleteMany(ctx context.Context) (int, error) {
	query := query.NewForceDelete()

	if err := fillTableCondition(ctx, query); err != nil {
		return 0, err
	}

	if err := build.Conditions(ctx, query); err != nil {
		return 0, err
	}

	logQuery(query)

	tx, err := Begin(ctx)
	if err != nil {
		return 0, err
	}

	_query, args, err := sqlx.Named(query.Query(), query.Arg())
	if err != nil {
		return 0, err
	}

	_query, args, err = sqlx.In(_query, args...)
	if err != nil {
		return 0, err
	}

	_query = tx.Rebind(_query)

	rows, err := tx.Exec(_query, args...)
	if err != nil {
		return 0, err
	}

	result, err := rows.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(result), nil
}
