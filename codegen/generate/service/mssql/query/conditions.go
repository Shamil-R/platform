package query

import (
	"fmt"
	"strings"
)

type conditions struct {
	query
	conditions []string
	and        []*conditions
	or         []*conditions
	not        []*conditions
}

func (q *conditions) block() string {
	conds := q.conditions

	if len(q.and) > 0 {
		var blocks []string
		for _, cond := range q.and {
			blocks = append(blocks, cond.block())
		}
		conds = append(conds, strings.Join(blocks, " OR"))
	}

	if len(q.or) > 0 {
		var blocks []string
		for _, cond := range q.or {
			blocks = append(blocks, cond.block())
		}
		conds = append(conds, strings.Join(blocks, " OR"))
	}

	if len(conds) == 0 {
		return ""
	}

	return "(" + strings.Join(conds, " AND") + ")"
}

func build(conds []string, op string) string {
	if len(conds) == 0 {
		return ""
	}
	return strings.Join(conds, " "+strings.TrimSpace(op))
}

func (q *conditions) AddСondition(column, op string, value interface{}) {
	if condFn, ok := conditionFuncs[op]; ok {
		placeholder := q.setArg(column, value)
		cond := condFn(column, placeholder)
		q.conditions = append(q.conditions, cond)
	}
}

func (q *conditions) And() Conditions {
	c := new(conditions)
	q.and = append(q.and, c)
	return c
}

func (q *conditions) Or() Conditions {
	c := new(conditions)
	q.or = append(q.or, c)
	return c
}

func (q *conditions) Not() Conditions {
	c := new(conditions)
	q.not = append(q.not, c)
	return c
}

func where(where string) string {
	if len(where) == 0 {
		return ""
	}
	return fmt.Sprintf("WHERE %s", where)
}

type conditionFunc func(column, placeholder string) string

func eq(column, placeholder string) string {
	return fmt.Sprintf("[%s] = :%s", column, placeholder)
}

func not(column, placeholder string) string {
	return fmt.Sprintf("[%s] != :%s", column, placeholder)
}

func lt(column, placeholder string) string {
	return fmt.Sprintf("[%s] < :%s", column, placeholder)
}

func lte(column, placeholder string) string {
	return fmt.Sprintf("[%s] <= :%s", column, placeholder)
}

func gt(column, placeholder string) string {
	return fmt.Sprintf("[%s] > :%s", column, placeholder)
}

func gte(column, placeholder string) string {
	return fmt.Sprintf("[%s] >= :%s", column, placeholder)
}

func contains(column, placeholder string) string {
	return fmt.Sprintf("[%s] LIKE CONCAT('%%', :%s, '%%')", column, placeholder)
}

func notContains(column, placeholder string) string {
	return fmt.Sprintf("[%s] NOT LIKE CONCAT('%%', :%s, '%%')", column, placeholder)
}

func startsWith(column, placeholder string) string {
	return fmt.Sprintf("[%s] LIKE CONCAT(:%s, '%%')", column, placeholder)
}

func notStartsWith(column, placeholder string) string {
	return fmt.Sprintf("[%s] NOT LIKE CONCAT(:%s, '%%')", column, placeholder)
}

func endsWith(column, placeholder string) string {
	return fmt.Sprintf("[%s] LIKE CONCAT('%%', :%s)", column, placeholder)
}

func notEndsWith(column, placeholder string) string {
	return fmt.Sprintf("[%s] NOT LIKE CONCAT('%%', :%s)", column, placeholder)
}

var conditionFuncs map[string]conditionFunc = map[string]conditionFunc{
	"eq":              eq,
	"not":             not,
	"lt":              lt,
	"lte":             lte,
	"gt":              gt,
	"gte":             gte,
	"contains":        contains,
	"not_contains":    notContains,
	"starts_with":     startsWith,
	"not_starts_with": notStartsWith,
	"ends_with":       endsWith,
	"not_ends_with":   notEndsWith,
}
