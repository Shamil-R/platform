package schema

import "github.com/vektah/gqlparser/ast"

const (
	ChildValueCreate  = "create"
	ChildValueConnect = "connect"
)

type ChildValue struct {
	*ast.ChildValue
	value *Value
}

func (v *ChildValue) Value() *Value {
	if v.value == nil {
		v.value = &Value{Value: v.ChildValue.Value}
	}
	return v.value
}

// func (v *ChildValue) Children() ChildValueList {
// 	return v.Value().Children()
// }

func (v *ChildValue) Directives() DirectiveList {
	return v.Value().Definition().Directives()
}

type ChildValueList []*ChildValue

func (l ChildValueList) Create() *ChildValue {
	return firstChildValue(l, isCreateChildValue)
}

func (l ChildValueList) Connect() *ChildValue {
	return firstChildValue(l, isConnectChildValue)
}

type childValueFilter func(childValue *ChildValue) bool

func isCreateChildValue(childValue *ChildValue) bool {
	return childValue.Name == ChildValueCreate
}

func isConnectChildValue(childValue *ChildValue) bool {
	return childValue.Name == ChildValueConnect
}

func hasChildValue(list ChildValueList, filter childValueFilter) bool {
	for _, childValue := range list {
		if filter(childValue) {
			return true
		}
	}
	return false
}

func firstChildValue(list ChildValueList, filter childValueFilter) *ChildValue {
	for _, childValue := range list {
		if filter(childValue) {
			return childValue
		}
	}
	return nil
}
