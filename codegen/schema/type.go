package schema

import "github.com/vektah/gqlparser/ast"

type Type struct {
	*ast.Type
	types DefinitionList
	elem  *Type
}

func (t *Type) IsDefinition() bool {
	return t.types.Contains(t)
}

func (t *Type) IsSlice() bool {
	return t.NamedType == "" && t.Elem() != nil
}

func (t *Type) Elem() *Type {
	if t.Type.Elem == nil {
		return nil
	}
	if t.elem == nil {
		t.elem = &Type{
			Type:  t.Type.Elem,
			types: t.types,
		}
	}
	return t.elem
}

func (t *Type) Name() string {
	if t.IsSlice() {
		return t.Elem().Name()
	}
	return t.Type.Name()
}
