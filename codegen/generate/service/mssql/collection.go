package mssql

import (
	"context"
	"fmt"
	"gitlab/nefco/platform/codegen/generate/service/mssql/query"

	"github.com/99designs/gqlgen/graphql"
)

func Collection(ctx context.Context, result interface{}) error {
	resCtx := graphql.GetResolverContext(ctx)

	query := query.Select(resCtx.Field.Field)

	fmt.Println(query.Query(), query.Arg())

	tx, err := Begin(ctx)
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareNamed(query.Query())
	if err != nil {
		return err
	}

	if err := stmt.Select(result, query.Arg()); err != nil {
		return err
	}

	return nil
}
