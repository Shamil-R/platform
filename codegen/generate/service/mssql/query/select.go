package query

import (
	"fmt"
	"strings"
)

type Select struct {
	conditions
	columns []string
}

func (q *Select) AddColumn(column, alias string) {
	col := fmt.Sprintf("[%s] AS %s", column, alias)
	q.columns = append(q.columns, col)
}

func (q *Select) Query() string {
	query := fmt.Sprintf(
		"SELECT %s FROM %s %s",
		strings.Join(q.columns, ", "),
		q.table,
		where(q.conditions.block()),
	)
	return query
}
