package build

import (
	"strings"

	"github.com/vektah/gqlparser/ast"
)

type Schema struct {
	*ast.Schema
}

func NewSchema(schema *ast.Schema) *Schema {
	return &Schema{schema}
}

func (s *Schema) Types() map[string]*ast.Definition {
	types := make(map[string]*ast.Definition)
	for key, def := range s.Schema.Types {
		isObject := def.Kind == ast.Object
		isEnum := def.Kind == ast.Enum
		if !strings.HasPrefix(def.Name, "__") && (isObject || isEnum) {
			types[key] = def
		}
	}
	return types
}
