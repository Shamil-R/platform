package query

import (
	"fmt"
	"strings"

	"github.com/vektah/gqlparser/ast"
)

func buildTable(input Input) string {
	var tableName string
	sel := input.Field().SelectionSet[0]
	switch sel := sel.(type) {
	case *ast.Field:
		def := sel.ObjectDefinition
		tableName = strings.ToLower(def.Name)
		if ok, val := directiveValue(def.Directives, "table", "name"); ok {
			tableName = val
		}
	}
	return fmt.Sprintf("[%s]", tableName)
}
