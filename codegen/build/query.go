package build

import (
	"strings"

	"github.com/jinzhu/inflection"

	"github.com/vektah/gqlparser/ast"
)

type ItemQuery struct {
	def *ast.Definition
}

func (q *ItemQuery) Name() string {
	return strings.ToLower(q.def.Name)
}

type ListQuery struct {
	def *ast.Definition
}

func (q *ListQuery) Name() string {
	return strings.ToLower(inflection.Plural(q.def.Name))
}
