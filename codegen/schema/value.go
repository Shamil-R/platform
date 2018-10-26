package schema

import "github.com/vektah/gqlparser/ast"

type Value struct {
	*ast.Value
}

func (v *Value) Children() ChildValueList {
	list := make(ChildValueList, 0, len(v.Value.Children))
	for _, child := range v.Value.Children {
		list = append(list, &ChildValue{child})
	}
	return list
}

func (v *Value) Definition() *Definition {
	// TODO: убрать Definition:
	return &Definition{Definition: v.Value.Definition}
}
