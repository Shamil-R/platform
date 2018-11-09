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

type Condition interface {
	AddСondition(column string, value interface{})
}

type Value interface {
	AddValue(column string, value interface{})
}

type query struct {
	table string
	arg   map[string]interface{}
}

func (q *query) addArg(key string, value interface{}) {
	if q.arg == nil {
		q.arg = map[string]interface{}{}
	}
	q.arg[key] = value
}

func (q *query) block() string {
	return q.table
}

func (q *query) SetTable(table string) {
	q.table = fmt.Sprintf("[%s]", table)
}

func (q *query) Arg() map[string]interface{} {
	return q.arg
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
	q.addArg(column, value)
}

/* type value struct {
	condition
	values []string
}

func (q *value) AddValue(column string, value interface{}) {
	q.values = append(q.values, column)
	q.addArg(column, value)
} */
