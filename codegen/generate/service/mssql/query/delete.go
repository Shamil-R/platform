package query

import (
	"fmt"

	"github.com/vektah/gqlparser/ast"
)

type deleteQuery struct {
	*query
	where Build
}

func (q deleteQuery) Query() string {
	return fmt.Sprintf(
		"DELETE FROM %s WHERE %s",
		buildTable(q),
		q.where.Build(q),
	)
}

func Delete(field *ast.Field) Query {
	return &deleteQuery{
		query: newQuery(field),
		where: &whereQuery{},
	}
}
