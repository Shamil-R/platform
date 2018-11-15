package query

import (
	"fmt"
)

type Delete struct {
	conditions
}

func (q *Delete) Query() string {
	query := fmt.Sprintf(
		"DELETE FROM %s %s",
		q.table,
		where(q.conditions.block()),
	)
	return query
}
