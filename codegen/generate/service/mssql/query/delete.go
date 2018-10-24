package query

import (
	"fmt"

	"github.com/vektah/gqlparser/ast"
)

type deleteQuery struct {
	*query
	table Build
	where Build
}

func (q deleteQuery) Query() string {
	return fmt.Sprintf(
		"DELETE FROM %s %s",
		q.table.Build(q),
		q.where.Build(q),
	)
}

func Delete(field *ast.Field) Query {
	return &deleteQuery{
		query: newQuery(field),
		table: &tableQuery{},
		where: &whereQuery{},
	}
}
