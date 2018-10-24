package mssql

import (
	"context"
	"fmt"
	"gitlab/nefco/platform/codegen/generate/service/mssql/query"

	"github.com/99designs/gqlgen/graphql"
)

func Item(ctx context.Context, result interface{}) error {
	resCtx := graphql.GetResolverContext(ctx)

	/* field := resCtx.Field

	for _, arg := range field.Arguments {
		fmt.Println(arg.Name)
		val := arg.Value
		def := val.Definition
		fmt.Println(def.Name)
		for _, child := range val.Children {
			fmt.Println(child.Name)
			childVal := child.Value
			fmt.Println(childVal.Definition.Name)
			f := def.Fields.ForName(child.Name)
			fmt.Println(f.Name)
			for _, d := range f.Directives {
				fmt.Println(d.Name)
			}
		}
	} */

	query := query.Select(resCtx.Field.Field)

	fmt.Println(query.Query())

	tx, err := Begin(ctx)
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareNamed(query.Query())
	if err != nil {
		return err
	}

	if err := stmt.Get(result, query.Arg()); err != nil {
		return err
	}

	return nil
}

func get(ctx context.Context, result interface{}) error {
	resCtx := graphql.GetResolverContext(ctx)

	tx, err := Begin(ctx)
	if err != nil {
		return err
	}

	query := query.Select(resCtx.Field.Field)

	fmt.Println("Query:", query.Query())

	stmt, err := tx.PrepareNamed(query.Query())
	if err != nil {
		return err
	}

	if err := stmt.Get(result, query.Arg()); err != nil {
		return err
	}

	return nil
}
