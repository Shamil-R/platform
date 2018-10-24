package query

import (
	"fmt"
	"strings"

	"github.com/vektah/gqlparser/ast"
)

type tableQuery struct{}

func (q tableQuery) Build(input Input) string {
	var table string
	sel := input.Field().SelectionSet[0]
	switch sel := sel.(type) {
	case *ast.Field:
		def := sel.ObjectDefinition
		table = strings.ToLower(def.Name)
		if ok, val := directiveValue(def.Directives, "table", "name"); ok {
			table = val
		}
	}
	return fmt.Sprintf("[%s]", table)
}
