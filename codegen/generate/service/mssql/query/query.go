package query

import (
	"fmt"
	"strings"
)

type Query interface {
	Query() string
	Arg() map[string]interface{}
}

type Table interface {
	SetTable(table string)
}

type Conditions interface {
	AddСondition(column string, value interface{})
}

type Values interface {
	AddValue(column string, value interface{})
}

type Columns interface {
	AddColumn(column string)
}

type query struct {
	arg   map[string]interface{}
	table string
}

func (q *query) setArg(key string, value interface{}) {
	if q.arg == nil {
		q.arg = map[string]interface{}{}
	}
	q.arg[key] = value
}

func (q *query) Arg() map[string]interface{} {
	return q.arg
}

func (q *query) SetTable(table string) {
	q.table = fmt.Sprintf("[%s]", table)
}

type condition struct {
	query
	conditions []string
}

func (q *condition) block() string {
	if len(q.conditions) == 0 {
		return ""
	}
	and := strings.Join(q.conditions, " AND")
	return fmt.Sprintf("WHERE %s", and)
}

func (q *condition) AddСondition(column string, value interface{}) {
	cond := fmt.Sprintf("[%s] = :%s", column, column)
	q.conditions = append(q.conditions, cond)
	q.setArg(column, value)
}
