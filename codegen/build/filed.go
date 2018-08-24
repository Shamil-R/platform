package build

import "github.com/vektah/gqlparser/ast"

type Field struct {
	field *ast.FieldDefinition
}

func (f *Field) Signature() string {
	return f.field.Type.Dump()
}
