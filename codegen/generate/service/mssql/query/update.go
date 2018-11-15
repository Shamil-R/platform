package query

import (
	"fmt"
	"strings"
)

type Update struct {
	conditions
	values []string
}

func (q *Update) AddValue(column string, value interface{}) {
	val := fmt.Sprintf("[%s] = :%s", column, column)
	q.values = append(q.values, val)
	q.setArg(column, value)
}

func (q *Update) Query() string {
	query := fmt.Sprintf(
		"UPDATE %s SET %s %s",
		q.table,
		strings.Join(q.values, ", "),
		where(q.conditions.block()),
	)
	return query
}
