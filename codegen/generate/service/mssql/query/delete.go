package query

import (
	"fmt"
)

type delete struct {
	*tableBlock
	*conditionsBlock
}

func NewDelete() *delete {
	return &delete{
		tableBlock:      newTableBlock(),
		conditionsBlock: newConditionsBlock(),
	}
}

func (q *delete) Query() string {
	query := fmt.Sprintf(
		"DELETE FROM %s %s",
		q.table,
		where(q.conditionsBlock.block()),
	)
	return query
}
