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
	overorderby := "order by %s %s"
	orderby := ""
	var field string
	var index string
	overfield := "(select null)"
	overindex := ""

	// Определяем столбец и направление сортировки
	if strings.HasSuffix(q.orderBy, "_ASC") {
		field = strings.TrimSuffix(q.orderBy, "_ASC")
		index = "ASC"
		overfield = field
		overindex = index
	} else if strings.HasSuffix(q.orderBy, "_DESC") {
		field = strings.TrimSuffix(q.orderBy, "_DESC")
		index = "DESC"
		overfield = field
		overindex = index
	}

	paginationCondition := fmt.Sprintf("and __num > %v", q.skip)

	if q.first > 0 {
		paginationCondition = fmt.Sprintf("%s and __num < %v", paginationCondition, q.skip + q.first + 1)
	} else if q.last > 0 {
		// при выводе last по умолчанию сортируют в обратном порядке, но если уже была определа сортировка,
		// то сортируем в противополжном ей направлению
		paginationCondition = fmt.Sprintf("%s and __num < %v", paginationCondition, q.skip + q.last + 1)
		overindex = "DESC"
		if overfield == "" {
			overfield = "id"
		}
		if index == "DESC" {
			overindex = "ASC"
		}
	}

	orderby = fmt.Sprintf("order by %s %s", field, index)
	overorderby = fmt.Sprintf("order by %s %s", overfield, overindex)

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
