package mssql

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
)

func Item(ctx context.Context, result interface{}) error {
	resCtx := graphql.GetResolverContext(ctx)

	field := resCtx.Field.Field

	arg := map[string]interface{}{}

	if argument := field.Arguments.ForName("where"); argument != nil {
		for _, child := range argument.Value.Children {
			arg[child.Name] = child.Value.Raw
		}
	}

	return itemByArg(ctx, result, arg)
}

func itemByArg(ctx context.Context, result interface{}, arg map[string]interface{}) error {
	resCtx := graphql.GetResolverContext(ctx)

	field := resCtx.Field.Field

	query := fmt.Sprintf(
		"SELECT \n\t%s\nFROM %s\n%s",
		querySelect(field),
		queryFrom(field),
		queryWhere(arg),
	)

	fmt.Println(query)

	tx, err := Begin(ctx)
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareNamed(query)
	if err != nil {
		return err
	}

	if err := stmt.Get(result, arg); err != nil {
		return err
	}

	return nil
}
