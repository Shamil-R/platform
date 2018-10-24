package query

import (
	"fmt"
	"strings"
)

type whereQuery map[string]interface{}

func (q whereQuery) Build(input Input) string {
	var conditions []string

	if ok, val := argumentValue(input, "where"); ok {
		conditions = make([]string, 0, len(val.Children))
		for _, child := range val.Children {
			input.Bind(child.Name, child.Value.Raw)
			columnName := columnName(val, child)
			condition := fmt.Sprintf("[%s] = :%s", columnName, child.Name)
			conditions = append(conditions, condition)
		}
	}

	for col, val := range q {
		input.Bind(col, val)
		condition := fmt.Sprintf("[%s] = :%s", col, col)
		conditions = append(conditions, condition)
	}

	if len(conditions) == 0 {
		return ""
	}

	return strings.Join(conditions, " AND")
}
