package mssql

import (
	"fmt"
	"strings"
)

type Table struct {
	table string
	arg   map[string]interface{}
}

func (q *Table) addArg(key string, value interface{}) {
	if q.arg == nil {
		q.arg = map[string]interface{}{}
	}
	q.arg[key] = value
}

func (q *Table) SetTable(table string) {
	q.table = fmt.Sprintf("[%s]", table)
}

func (q *Table) Query() string {
	return q.table
}

func (q *Table) Arg() map[string]interface{} {
	return q.arg
}

type Condition struct {
	Table
	conditions []string
}

func (q *Condition) AddСondition(column string, value interface{}) {
	cond := fmt.Sprintf("[%s] = :%s", column, column)
	q.conditions = append(q.conditions, cond)
	q.addArg(column, value)
}

func (q *Condition) Query() string {
	if len(q.conditions) == 0 {
		return ""
	}
	and := strings.Join(q.conditions, " AND")
	return fmt.Sprintf("WHERE %s", and)
}

type Value struct {
	Condition
	values []string
}

func (q *Value) AddValue(column string, value interface{}) {
	q.values = append(q.values, column)
	q.addArg(column, value)
}

/* func collectPagination(field *ast.Field) ([]string, map[string]interface{}) {
	var query []string
	arg := make(map[string]interface{})

	args := paginationArgs(field)
	if len(args) == 0 {
		return query, arg
	}

	skip := args["skip"]

	if skip != nil {
		cond := "[num] > :skip"
		query = append(query, cond)
		arg["skip"] = skip.Value.Raw
	}

	return query, arg
}

func paginationArgs(field *ast.Field) map[string]*ast.Argument {
	args := make(map[string]*ast.Argument)
	names := []string{"skip", "first", "last"}
	for _, name := range names {
		if arg := field.Arguments.ForName(name); arg != nil {
			args[name] = arg
		}
	}
	return args
} */
