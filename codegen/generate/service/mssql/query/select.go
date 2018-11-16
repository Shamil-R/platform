package query

import (
	"fmt"
	"strings"
)

type zelect struct {
	*tableBlock
	*conditionsBlock
	columns []string
}

func NewSelect() *zelect {
	return &zelect{
		tableBlock:      newTableBlock(),
		conditionsBlock: newConditionsBlock(),
	}
}

func (q *zelect) AddColumn(column, alias string) {
	col := fmt.Sprintf("[%s] AS %s", column, alias)
	q.columns = append(q.columns, col)
}

func (q *zelect) Query() string {
	query := fmt.Sprintf(
		"SELECT %s FROM %s %s",
		strings.Join(q.columns, ", "),
		q.table,
		where(q.conditionsBlock.block()),
	)
	return query
}
