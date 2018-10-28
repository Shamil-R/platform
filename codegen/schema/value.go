package schema

import "github.com/vektah/gqlparser/ast"

type Value struct {
	*ast.Value
	children ChildValueList
}

func (v *Value) Children() ChildValueList {
	if v.children != nil {
		return v.children
	}
	v.children = make(ChildValueList, 0, len(v.Value.Children))
	for _, child := range v.Value.Children {
		v.children = append(v.children, &ChildValue{ChildValue: child})
	}
	return v.children
}

func (v *Value) Definition() *Definition {
	// TODO: убрать Definition:
	return &Definition{Definition: v.Value.Definition}
}
