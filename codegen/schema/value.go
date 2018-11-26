package schema

import (
	"github.com/vektah/gqlparser/ast"
	"strconv"
)

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

func (v *Value) Conv() interface{} {
	switch v.Kind {
	case ast.IntValue:
		n, _ := strconv.Atoi(v.Raw)
		return n
	case ast.ListValue:
		var arr []string
		for _, child := range v.Value.Children {
			arr = append(arr, child.Value.Raw)
		}
		return arr
	}
	return v.Raw
}
