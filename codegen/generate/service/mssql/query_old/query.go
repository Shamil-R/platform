package query_old

import (
	"github.com/vektah/gqlparser/ast"
)

type Bind interface {
	Bind(placeholder string, value interface{})
}

type Input interface {
	Bind
	Field() *ast.Field
}

type Build interface {
	Build(input Input) string
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

func newQuery(field *ast.Field) *query {
	return &query{field, make(map[string]interface{})}
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
