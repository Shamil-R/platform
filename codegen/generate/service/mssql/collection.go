package mssql

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
)

func Collection(ctx context.Context, result interface{}) error {
	resCtx := graphql.GetResolverContext(ctx)

	field := resCtx.Field.Field

	query, arg := qSelect(field)

	// arg := map[string]interface{}{}

	// query := fmt.Sprintf(
	// 	"SELECT \n\t%s\nFROM %s",
	// 	querySelect(field),
	// 	queryFrom(field),
	// )

	fmt.Println(query)

	tx, err := Begin(ctx)
	if err != nil {
		return err
	}

	if len(arg) == 0 {
		if err := tx.Select(result, query); err != nil {
			return err
		}
	} else {
		stmt, err := tx.PrepareNamed(query)
		if err != nil {
			return err
		}

		if err := stmt.Select(result, arg); err != nil {
			return err
		}
	}

	return nil
}
