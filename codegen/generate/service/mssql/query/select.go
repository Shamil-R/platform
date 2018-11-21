package query

import (
	"fmt"
	"strings"
)

type zelect struct {
	*tableBlock
	*conditionsBlock
	columns []string
	skip int
	first int
	last int
	orderBy string
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
	var field string
	var index string
	if strings.HasSuffix(q.orderBy, "_ASC") {
		field = strings.TrimSuffix(q.orderBy, "_ASC")
		index = "ASC"
	} else if strings.HasSuffix(q.orderBy, "_DESC") {
		field = strings.TrimSuffix(q.orderBy, "_DESC")
		index = "DESC"
	}

	overorderby := "order by %s %s"
	orderby := ""
	paginationCondition := fmt.Sprintf("and __num > %v", q.skip)

	if q.first > 0 {
		paginationCondition = fmt.Sprintf("%s and __num < %v", paginationCondition, q.skip + q.first + 1)
	} else if q.last > 0 {
		paginationCondition = fmt.Sprintf("%s and __num < %v", paginationCondition, q.skip + q.last + 1)
		if field == "" {
			field = "id"
		}
	}

	if field != "" {
		orderby = fmt.Sprintf(overorderby, field, index)
		overorderby = fmt.Sprintf(overorderby, field, "DESC")
	} else {
		overorderby = fmt.Sprintf(overorderby, "(select null)", "")
	}

	query := fmt.Sprintf(
		"SELECT %s from (SELECT ROW_NUMBER() over (%s) as __num, %s FROM %s %s ) a where 1=1 %s %s",
		strings.Join(q.columns, ", "),
		overorderby,
		strings.Join(q.columns, ", "),
		q.table,
		where(q.conditionsBlock.block()),
		paginationCondition,
		orderby,
	)
	return query
}

func (q *zelect) SetSkip(skip int) {
	q.skip = skip
}

func (q *zelect) SetFirst(first int) {
	q.first = first
}

func (q *zelect) SetLast(last int) {
	q.last = last
}

func (q *zelect) SetOrder(orderBy string) {
	q.orderBy = orderBy
}
