package query

import (
	"fmt"
	"strconv"
	"strings"
)

func buildPagination(input Input) string {
	isSkip, skip := argumentValue(input, "skip")
	isFirst, first := argumentValue(input, "first")
	isLast, last := argumentValue(input, "last")
	if !isSkip && !isFirst && !isLast {
		return ""
	}

	subWhere := whereQuery{}.Build(input)

	if len(subWhere) != 0 {
		subWhere = fmt.Sprintf("WHERE %s", subWhere)
	}

	order := "ASC"

	var where []string

	if isSkip {
		v, _ := strconv.Atoi(skip.Raw)
		input.Bind("skip", v)
		where = append(where, "[num] > :skip")
	} else {
		input.Bind("skip", 0)
	}

	if isFirst {
		v, _ := strconv.Atoi(first.Raw)
		input.Bind("first", v)
		where = append(where, "[num] <= :skip + :first")
	} else if isLast {
		v, _ := strconv.Atoi(last.Raw)
		input.Bind("last", v)
		where = append(where, "[num] <= :skip + :last")
		order = "DESC"
	}

	column := fmt.Sprintf("ROW_NUMBER() OVER (ORDER BY [id] %s) AS num", order)

	return fmt.Sprintf(
		"WITH Ordered AS (SELECT %s FROM %s %s) SELECT %s FROM Ordered WHERE %s ORDER BY [id] ASC",
		selectColumnsQuery([]string{column}).Build(input),
		tableQuery{}.Build(input),
		subWhere,
		selectColumnsQuery{}.Build(input),
		strings.Join(where, " AND"),
	)
}
