package mssql

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"gitlab/nefco/platform/codegen/generate/service/mssql/build"
	"gitlab/nefco/platform/codegen/generate/service/mssql/query"
)

func Collection(ctx context.Context, result interface{}) error {
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

	if err := build.Order(ctx, query); err != nil {
		return err
	}

	logQuery(query)

	tx, err := Begin(ctx)
	if err != nil {
		return err
	}

	/*stmt, err := tx.PrepareNamed(query.Query())
	if err != nil {
		return err
	}*/

	_query, args, err := sqlx.Named(query.Query(), query.Arg())
	if err != nil {
		return err
	}
	_query, args, err = sqlx.In(_query, args...)
	if err != nil {
		return err
	}

	rows, err := tx.NamedQuery(query.Query(), query.Arg())
	if err != nil {
		return err
	}
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		fmt.Println("id=",id)
	}
	_query = tx.Rebind(_query)
	fmt.Println("--------------")
	fmt.Println(_query)

	stmt, err := tx.PrepareNamed(_query)
	if err != nil {
		return err
	}
	if err := stmt.Select(result, args); err != nil {
		return err
	}

	return nil
}
