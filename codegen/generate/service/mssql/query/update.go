package query

import (
	"fmt"
	"strings"

	"github.com/vektah/gqlparser/ast"
)

type updateQuery struct {
	*query
	table   Build
	columns Build
	where   Build
}

func (q updateQuery) Query() string {
	return fmt.Sprintf(
		"UPDATE %s SET %s %s",
		q.table.Build(q),
		q.columns.Build(q),
		q.where.Build(q),
	)
}

func Update(field *ast.Field) Query {
	return &updateQuery{
		query:   newQuery(field),
		table:   &tableQuery{},
		columns: &updateColumnsQuery{},
		where:   &whereQuery{},
	}
}

type updateColumnsQuery struct{}

func (q updateColumnsQuery) Build(input Input) string {
	ok, val := argumentValue(input, "data")
	if !ok {
		return ""
	}
	columns := make([]string, 0, len(val.Children))
	for _, child := range val.Children {
		input.Bind(child.Name, child.Value.Raw)
		columnName := columnName(val, child)
		column := fmt.Sprintf("[%s] = :%s", columnName, child.Name)
		columns = append(columns, column)
	}
	return strings.Join(columns, ", ")
}
