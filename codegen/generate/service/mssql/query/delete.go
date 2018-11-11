package query

import (
	"fmt"
)

type Delete struct {
	condition
}

func (q *Delete) Query() string {
	query := fmt.Sprintf(
		"DELETE FROM %s %s",
		q.table,
		q.condition.block(),
	)
	return query
}
