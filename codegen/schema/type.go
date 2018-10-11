package schema

import "github.com/vektah/gqlparser/ast"

type Type struct {
	*ast.Type
	schema *Schema
}

func (t *Type) IsDefinition() bool {
	return t.schema.Types().Objects().Contains(t)
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

func (t *Type) Name() string {
	if t.IsSlice() {
		return t.Elem().Name()
	}
	return t.Type.Name()
}

func (t *Type) Definition() *Definition {
	return t.schema.Types().ByName(t.Name())
}
