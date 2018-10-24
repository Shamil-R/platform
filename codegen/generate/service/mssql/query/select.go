package query

import (
	"fmt"
	"strings"

	"github.com/vektah/gqlparser/ast"
)

type selectQuery struct {
	*query
	columns Build
	table   Build
	where   Build
}

func (q selectQuery) Query() string {
	return fmt.Sprintf(
		"SELECT %s FROM %s %s",
		q.columns.Build(q),
		q.table.Build(q),
		q.where.Build(q),
	)
}

func Select(field *ast.Field) Query {
	return &selectQuery{
		query:   newQuery(field),
		columns: &selectColumnsQuery{},
		table:   &tableQuery{},
		where:   &whereQuery{},
	}
}

func SelectWhere(field *ast.Field, where map[string]interface{}) Query {
	wq := whereQuery(where)
	return &selectQuery{
		query:   newQuery(field),
		columns: &selectColumnsQuery{},
		table:   &tableQuery{},
		where:   &wq,
	}
}

type selectColumnsQuery struct{}

func (q selectColumnsQuery) Build(input Input) string {
	set := input.Field().SelectionSet
	columns := make([]string, 0, len(set))
	for _, sel := range set {
		switch sel := sel.(type) {
		case *ast.Field:
			column := fmt.Sprintf("[%s]", sel.Name)
			def := sel.ObjectDefinition
			if ok, val := directiveValue(def.Directives, "field", "name"); ok {
				column = fmt.Sprintf("%s AS %s", column, val)
			}
			columns = append(columns, column)
		}
	}
	return strings.Join(columns, ", ")
}
