package query

import (
	"fmt"
)

type restore struct {
	*tableBlock
	*conditionsBlock
}

func NewRestore() *restore {
	return &restore{
		tableBlock:      newTableBlock(),
		conditionsBlock: newConditionsBlock(),
	}
}

func (q *restore) Query() string {
	query := fmt.Sprintf(
		"DELETE FROM %s %s",
		q.table,
		where(q.conditionsBlock.block()),
	)
	return query
}
