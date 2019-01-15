package query_old

import (
	"fmt"
	"strings"

	"github.com/vektah/gqlparser/ast"
)

type insertQuery struct {
	*query
	columns Build
	values  Build
}

func (q insertQuery) Query() string {
	return fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		buildTable(q),
		q.columns.Build(q),
		q.values.Build(q),
	)
}

func Insert(field *ast.Field) Query {
	return &insertQuery{
		query:   newQuery(field),
		columns: &insertColumnsQuery{},
		values:  &valuesQuery{},
	}
}

type insertColumnsQuery struct{}

func (q insertColumnsQuery) Build(input Input) string {
	ok, val := argumentValue(input, "data")
	if !ok {
		return ""
	}
	columns := make([]string, 0, len(val.Children))
	for _, child := range val.Children {
		input.Bind(child.Name, child.Value.Raw)
		columnName := columnName(val, child)
		column := fmt.Sprintf("[%s]", columnName)
		columns = append(columns, column)
	}
	return strings.Join(columns, ", ")
}

type valuesQuery struct{}

func (q valuesQuery) Build(input Input) string {
	ok, val := argumentValue(input, "data")
	if !ok {
		return ""
	}
	values := make([]string, 0, len(val.Children))
	for _, child := range val.Children {
		var (
			placeholder string
			value       *ast.Value
		)
		if connect := child.Value.Children.ForName("connect"); connect != nil {
			placeholder = "id"
			value = connect.Children.ForName(placeholder)
		} else {
			placeholder = child.Name
			value = child.Value
		}
		input.Bind(placeholder, value.Raw)
		val := fmt.Sprintf(":%s", placeholder)
		values = append(values, val)
	}
	return strings.Join(values, ", ")
}
