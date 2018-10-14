package mssql

import (
	"context"
	"fmt"
	"strings"

	"github.com/99designs/gqlgen/graphql"
)

func Create(ctx context.Context, result interface{}) error {
	resCtx := graphql.GetResolverContext(ctx)

	field := resCtx.Field.Field

	// fmt.Println("field:", field.Name, field.Alias, field.Definition.Type)

	// for _, arg := range field.Arguments {
	// 	fmt.Println("arg:", arg.Name, arg.Value)

	// 	def := arg.Value.Definition
	// 	if def != nil {
	// 		fmt.Println("def:", def.Name)
	// 	}
	// }

	// for _, sel := range resCtx.Field.Selections {
	// 	switch sel := sel.(type) {
	// 	case *ast.Field:
	// 		fmt.Println("sel:", sel.Name)
	// 	}
	// }

	var queryInsertColumns []string
	var queryInsertValues []string
	arg := map[string]interface{}{}

	if argument := field.Arguments.ForName("data"); argument != nil {
		for _, child := range argument.Value.Children {
			queryInsertColumns = append(queryInsertColumns, fmt.Sprintf("[%s]", child.Name))
			queryInsertValues = append(queryInsertValues, fmt.Sprintf(":%s", child.Name))
			arg[child.Name] = child.Value.Raw
		}
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (\n%s\n) VALUES (\n%s\n)",
		queryFrom(field),
		strings.Join(queryInsertColumns, ",\n\t"),
		strings.Join(queryInsertValues, ",\n\t"),
	)

	fmt.Println(query)

	tx, err := Begin(ctx)
	if err != nil {
		return err
	}

	res, err := tx.NamedExec(query, arg)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	fmt.Println(id)

	arg = map[string]interface{}{
		"id": id,
	}

	return itemByArg(ctx, result, arg)
}
