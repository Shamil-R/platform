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

type Pagination interface {
	SetSkip(int)
	SetFirst(int)
	SetLast(int)
}

type Sort interface {
	SetOrder(string, string)
}

type Values interface {
	AddValue(column string, value interface{})
}

type Columns interface {
	AddColumn(column, alias string)
}

type arg struct {
	m map[string]interface{}
}

func newArg() *arg {
	return &arg{m: map[string]interface{}{}}
}

func (a *arg) setArg(key string, value interface{}) string {
	l := len(a.m)
	k := key + "_" + strconv.Itoa(l)
	a.m[k] = value
	return k
}

func (a *arg) Arg() map[string]interface{} {
	return a.m
}

type tableBlock struct {
	table string
}

func newTableBlock() *tableBlock {
	return &tableBlock{}
}

func (t *tableBlock) SetTable(table string) {
	t.table = fmt.Sprintf("[%s]", table)
}
