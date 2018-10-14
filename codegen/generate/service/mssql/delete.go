package mssql

import (
	"context"
	"fmt"
	"strings"

	"github.com/99designs/gqlgen/graphql"
)

func Delete(ctx context.Context, result interface{}) error {

	resCtx := graphql.GetResolverContext(ctx)

	field := resCtx.Field.Field

	var queryWhere []string
	arg := map[string]interface{}{}

	if argument := field.Arguments.ForName("where"); argument != nil {
		for _, child := range argument.Value.Children {
			col := fmt.Sprintf("[%s] = :%s", child.Name, child.Name)
			queryWhere = append(queryWhere, col)
			arg[child.Name] = child.Value.Raw
		}
	}

	query := fmt.Sprintf(
		"DELETE FROM %s\nWHERE\n%s",
		queryFrom(field),
		strings.Join(queryWhere, ",\n\t"),
	)

	fmt.Println(query)

	tx, err := Begin(ctx)
	if err != nil {
		return err
	}

	if err := itemByArg(ctx, result, arg); err != nil {
		return err
	}

	if _, err := tx.NamedExec(query, arg); err != nil {
		return err
	}

	return nil
}
