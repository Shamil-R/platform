package query

import (
	"fmt"
	"strings"
)

type update struct {
	*tableBlock
	*conditionsBlock
	values []string
}

func NewUpdate() *update {
	return &update{
		tableBlock:      newTableBlock(),
		conditionsBlock: newConditionsBlock(),
	}
}

func (q *update) AddValue(column string, value interface{}) {
	arg := q.setArg(column, value)
	val := fmt.Sprintf("[%s] = :%s", column, arg)
	q.values = append(q.values, val)

}

func (q *update) Query() string {
	query := fmt.Sprintf(
		"UPDATE %s SET %s %s",
		q.table,
		strings.Join(q.values, ", "),
		where(q.conditionsBlock.block()),
	)
	return query
}
