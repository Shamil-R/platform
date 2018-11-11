package query

import (
	"fmt"
	"strings"
)

type Insert struct {
	query
	columns []string
	values  []string
}

func (q *Insert) AddValue(column string, value interface{}) {
	col := fmt.Sprintf("[%s]", column)
	q.columns = append(q.columns, col)
	val := fmt.Sprintf(":%s", column)
	q.values = append(q.values, val)
	q.setArg(column, value)
}

func (q *Insert) Query() string {
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		q.table,
		strings.Join(q.columns, ", "),
		strings.Join(q.values, ", "),
	)
	return query
}
