package query

import (
	"fmt"
	"strings"

	"github.com/vektah/gqlparser/ast"
)

type selectQuery struct {
	*query
	whereQuery
}

func (q selectQuery) Query() string {
	query := buildPagination(q)
	if len(query) != 0 {
		return query
	}

	where := q.whereQuery.Build(q)
	if len(where) != 0 {
		where = fmt.Sprintf("WHERE %s", where)
	}

	return fmt.Sprintf(
		"SELECT %s FROM %s %s",
		selectColumnsQuery{}.Build(q),
		tableQuery{}.Build(q),
		where,
	)
}

func Select(field *ast.Field) Query {
	return &selectQuery{
		query: newQuery(field),
	}
}

func SelectWhere(field *ast.Field, where map[string]interface{}) Query {
	return &selectQuery{
		query:      newQuery(field),
		whereQuery: whereQuery(where),
	}
}

type selectColumnsQuery []string

func newSelectColumnsQuery(columns []string) *selectColumnsQuery {
	q := selectColumnsQuery(columns)
	return &q
}

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

	for _, column := range q {
		columns = append(columns, column)
	}

	if len(columns) == 0 {
		return "*"
	}

	return strings.Join(columns, ", ")
}
