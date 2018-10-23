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

func Select(field *ast.Field, where map[string]interface{}) Query {
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

type tableQuery struct{}

func (q tableQuery) Build(input Input) string {
	var table string
	sel := input.Field().SelectionSet[0]
	switch sel := sel.(type) {
	case *ast.Field:
		def := sel.ObjectDefinition
		table = strings.ToLower(def.Name)
		if ok, val := directiveValue(def.Directives, "table", "name"); ok {
			table = val
		}
	}
	return fmt.Sprintf("[%s]", table)
}

type whereQuery map[string]interface{}

func (q whereQuery) Build(input Input) string {
	var conditions []string
	ok, children := argumentChildren(input, "where")
	if ok {
		conditions = make([]string, 0, len(children))
		for _, child := range children {
			input.Bind(child.Name, child.Value.Raw)
			condition := fmt.Sprintf("[%s] = :%s", child.Name, child.Name)
			conditions = append(conditions, condition)
		}
	}
	for col, val := range q {
		input.Bind(col, val)
		condition := fmt.Sprintf("[%s] = :%s", col, col)
		conditions = append(conditions, condition)
	}
	if len(conditions) == 0 {
		return ""
	}
	return fmt.Sprintf("WHERE %s", strings.Join(conditions, " AND"))
}
