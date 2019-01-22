package query

import (
	"fmt"
)

func order(orderField, orderIndex string) string {
	if len(orderField) == 0 {
		return ""
	}
	if len(orderIndex) == 0 {
		orderIndex = "ASC"
	}
	return fmt.Sprintf("ORDER BY %s %s", orderField, orderIndex)
}
