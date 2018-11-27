package query

import (
	"fmt"
	"strings"
)

type insert struct {
	*arg
	*tableBlock
	columns []string
	values  []string
}

func NewInsert() *insert {
	return &insert{
		arg:        newArg(),
		tableBlock: newTableBlock(),
	}
}

func (q *insert) AddValue(column string, value interface{}) {
	col := fmt.Sprintf("[%s]", column)
	q.columns = append(q.columns, col)
	val := q.setArg(column, column)
	q.values = append(q.values, fmt.Sprintf(":%s", val))

}

func (q *insert) Query() string {
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		q.table,
		strings.Join(q.columns, ", "),
		strings.Join(q.values, ", "),
	)
	return query
}
