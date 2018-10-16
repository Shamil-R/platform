package query

import "go/ast"

type Bind interface {
	Bind(placeholder string, value interface{})
}

type Field interface {
	Bind
	Field() *ast.Field
}

type Build interface {
	Build(field Field) string
}

type Arg interface {
	Arg() map[string]interface{}
}

type Query interface {
	Arg
	Query() string
}

type query struct {
	field *ast.Field
	arg   map[string]interface{}
}

func (q *query) Bind(placeholder string, value interface{}) {
	q.arg[placeholder] = value
}

func (q *query) Arg() map[string]interface{} {
	return q.arg
}

func (q *query) Field() *ast.Field {
	return q.field
}
