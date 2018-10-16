package query

import (
	"fmt"
	"go/ast"
)

type selectQuery struct {
	*query
	columns   Build
	table     Build
	condition Build
}

func (q selectQuery) Query() string {
	return fmt.Sprintf(
		"SELECT %s %s %s",
		q.columns.Build(q),
		q.table.Build(q),
		q.condition.Build(q),
	)
}

func Select(field *ast.Field) Query {
	return &selectQuery{
		query: &query{
			field: field,
			arg:   make(map[string]interface{}),
		},
		condition: &whereQuery{},
	}
}

type whereQuery struct {
	*query
}
