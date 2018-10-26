package schema

import "github.com/vektah/gqlparser/ast"

type ChildValue struct {
	*ast.ChildValue
}

func (v *ChildValue) Value() *Value {
	return &Value{v.ChildValue.Value}
}

type ChildValueList []*ChildValue
