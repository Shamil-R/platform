package query

import (
	"fmt"
	"strings"
)

type Update struct {
	condition
	values []string
}

func (q *Update) block() string {
	values := make([]string, 0, len(q.values))
	for _, col := range q.values {
		val := fmt.Sprintf("[%s] = :%s", col, col)
		values = append(values, val)
	}
	return strings.Join(values, ", ")
}

func (q *Update) AddValue(column string, value interface{}) {
	q.values = append(q.values, column)
	q.addArg(column, value)
}

func (q *Update) Query() string {
	query := fmt.Sprintf(
		"UPDATE %s SET %s %s",
		q.query.block(),
		q.block(),
		q.condition.block(),
	)
	return query
}
