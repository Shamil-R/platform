package mssql

import (
	"context"
	"fmt"
	"gitlab/nefco/platform/codegen/generate/service/mssql/query"

	"github.com/99designs/gqlgen/graphql"
)

func Delete(ctx context.Context, result interface{}) error {
	resCtx := graphql.GetResolverContext(ctx)

	tx, err := Begin(ctx)
	if err != nil {
		return err
	}

	querySelect := query.Select(resCtx.Field.Field)

	fmt.Println(querySelect.Query())

	stmt, err := tx.PrepareNamed(querySelect.Query())
	if err != nil {
		return err
	}

	if err := stmt.Get(result, querySelect.Arg()); err != nil {
		return err
	}

	queryDelete := query.Delete(resCtx.Field.Field)

	fmt.Println(queryDelete.Query())

	if _, err := tx.NamedExec(queryDelete.Query(), queryDelete.Arg()); err != nil {
		return err
	}

	return nil
}
