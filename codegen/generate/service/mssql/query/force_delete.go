package query

import (
	"fmt"
)

type forceDelete struct {
	*tableBlock
	*conditionsBlock
}

func NewForceDelete() *forceDelete {
	return &forceDelete{
		tableBlock:      newTableBlock(),
		conditionsBlock: newConditionsBlock(),
	}
}

func (q *forceDelete) Query() string {
	query := fmt.Sprintf(
		"DELETE FROM %s %s",
		q.table,
		where(q.conditionsBlock.block()),
	)
	return query
}
