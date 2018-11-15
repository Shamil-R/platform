package query

import (
	"fmt"
	"strconv"
)

type Query interface {
	Query() string
	Arg() map[string]interface{}
}

type Table interface {
	SetTable(table string)
}

type Conditions interface {
	AddСondition(column, condition string, value interface{})
	And() Conditions
	Or() Conditions
	Not() Conditions
}

type Values interface {
	AddValue(column string, value interface{})
}

type Columns interface {
	AddColumn(column, alias string)
}

type query struct {
	arg   map[string]interface{}
	table string
}

func (q *query) setArg(key string, value interface{}) string {
	if q.arg == nil {
		q.arg = map[string]interface{}{}
	}
	k := key + "_" + strconv.Itoa(len(q.arg))
	q.arg[k] = value
	return k
}

func (q *query) Arg() map[string]interface{} {
	return q.arg
}

func (q *query) SetTable(table string) {
	q.table = fmt.Sprintf("[%s]", table)
}
