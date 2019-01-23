package mssql

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"gitlab/nefco/platform/codegen/generate/service/mssql/build"
	"gitlab/nefco/platform/codegen/generate/service/mssql/query"
)

func RestoreMany(ctx context.Context, result interface{}) (error) {
	query := query.NewUpdate()

	if err := fillTableCondition(ctx, query); err != nil {
		return err
	}

	dirName := "softDelete"
	argName := "deleteField"
	fieldName, err := getDefaultValues(ctx, dirName, argName)
	if err != nil {
		return err
	}
	query.AddValue(fieldName, sql.NullString{})

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

	rows, err := tx.Exec(_query, args...)
	if err != nil {
		return err
	}

	_, err = rows.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}
