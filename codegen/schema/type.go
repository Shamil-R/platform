package schema

import "github.com/vektah/gqlparser/ast"

type Type struct {
	*ast.Type
	schema *Schema
}

func (t *Type) IsObject() bool {
	name := t.NamedType
	// if t.IsSlice() {
	// 	name = t.Elem().NamedType
	// }
	return t.schema.Types().Objects().ByName(name) != nil
}

func (t *Type) IsSlice() bool {
	return t.NamedType == "" && t.Elem() != nil
}

func (t *Type) Elem() *Type {
	if t.Type.Elem == nil {
		return nil
	}
	return &Type{t.Type.Elem, t.schema}
}
