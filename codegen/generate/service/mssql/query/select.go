package query

import (
	"fmt"
	"strings"
)

type Select struct {
	condition
	columns []string
}

func (q *Select) AddColumn(column string) {
	col := fmt.Sprintf("[%s]", column)
	q.columns = append(q.columns, col)
}

func (q *Select) Query() string {
	query := fmt.Sprintf(
		"SELECT %s FROM %s %s",
		strings.Join(q.columns, ", "),
		q.table,
		q.condition.block(),
	)
	return query
}
