package query

import (
	"fmt"
	"strings"
)

type zelect struct {
	*tableBlock
	*conditionsBlock
	columns []string
	aliases []string
	skip int
	first int
	last int
	orderField string
	orderIndex string
	withTrashed bool
	onlyTrashed bool
	trashedFieldName string
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
	q.aliases = append(q.aliases, alias)
}

func (q *zelect) Query() string {
	overorderby := "order by %s %s"
	orderby := ""
	overfield := "(select null)"
	overindex := ""

	// Определяем столбец и направление сортировки
	if q.orderIndex == "ASC" {
		overfield = q.orderField
		overindex = q.orderIndex
	} else if q.orderIndex == "DESC" {
		overfield = q.orderField
		overindex = q.orderIndex
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
		if q.orderIndex == "DESC" {
			overindex = "ASC"
		}
	}

	if q.orderField != "" {
		orderby = fmt.Sprintf("order by %s %s", q.orderField, q.orderIndex)
	}
	overorderby = fmt.Sprintf("order by %s %s", overfield, overindex)


	if q.trashedFieldName != "" {
		if q.onlyTrashed {
			q.AddСondition(q.trashedFieldName, "is_not", "null")
		} else if q.withTrashed {
		} else {
			q.AddСondition(q.trashedFieldName, "is", "null")
		}
	}


	query := fmt.Sprintf(
		"SELECT %s from (SELECT ROW_NUMBER() over (%s) as __num, %s FROM %s %s ) a where 1=1  %s %s",
		strings.Join(q.aliases, ", "),
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

func (q *zelect) SetOrder(orderField string, orderIndex string) {
	q.orderField = orderField
	q.orderIndex = orderIndex
}

func (q *zelect) SetTrashed(withTrashed bool, onlyTrashed bool) {
	q.withTrashed = withTrashed
	q.onlyTrashed = onlyTrashed
}

func (q *zelect) SetTrashedFieldName(column string) {
	q.trashedFieldName = fmt.Sprintf("%s", column)
}
