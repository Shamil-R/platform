package query_old

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

	numberOrder := "ASC"

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
		numberOrder = "DESC"
	}

	numberColumn := fmt.Sprintf(
		"ROW_NUMBER() OVER (ORDER BY [id] %s) AS num",
		numberOrder,
	)

	numberWhere := whereQuery{}.Build(input)

	if len(numberWhere) != 0 {
		numberWhere = fmt.Sprintf("WHERE %s", numberWhere)
	}

	numberQuery := fmt.Sprintf(
		"WITH Ordered AS (SELECT %s, %s FROM %s %s)",
		selectColumnsQuery{}.Build(input),
		numberColumn,
		buildTable(input),
		numberWhere,
	)

	return fmt.Sprintf(
		"%s SELECT %s FROM Ordered WHERE %s ORDER BY [id] ASC",
		numberQuery,
		selectColumnsQuery{}.Build(input),
		strings.Join(where, " AND"),
	)
}
