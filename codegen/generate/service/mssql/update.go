package mssql

import (
	"context"
	"fmt"
	"strings"

	"github.com/99designs/gqlgen/graphql"
)

func Update(ctx context.Context, result interface{}) error {
	resCtx := graphql.GetResolverContext(ctx)

	field := resCtx.Field.Field

	var queryUpdate []string
	var queryWhere []string
	arg := map[string]interface{}{}
	argWhere := map[string]interface{}{}

	if argument := field.Arguments.ForName("data"); argument != nil {
		for _, child := range argument.Value.Children {
			col := fmt.Sprintf("[%s] = :%s", child.Name, child.Name)
			queryUpdate = append(queryUpdate, col)
			arg[child.Name] = child.Value.Raw
		}
	}

	if argument := field.Arguments.ForName("where"); argument != nil {
		for _, child := range argument.Value.Children {
			col := fmt.Sprintf("[%s] = :%s", child.Name, child.Name)
			queryWhere = append(queryWhere, col)
			arg[child.Name] = child.Value.Raw
			argWhere[child.Name] = child.Value.Raw
		}
	}

	query := fmt.Sprintf(
		"UPDATE %s SET\n%s\nWHERE\n%s",
		queryFrom(field),
		strings.Join(queryUpdate, ",\n\t"),
		strings.Join(queryWhere, ",\n\t"),
	)

	fmt.Println(query)

	tx, err := Begin(ctx)
	if err != nil {
		return err
	}

	if _, err := tx.NamedExec(query, arg); err != nil {
		return err
	}

	return itemByArg(ctx, result, argWhere)
}
