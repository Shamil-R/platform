package schema

import "github.com/vektah/gqlparser/ast"

type ArgumentDefinition struct {
	*ast.ArgumentDefinition
	schema *Schema
}

func (a *ArgumentDefinition) Type() *Type {
	return &Type{a.ArgumentDefinition.Type, a.schema}
}

type ArgumentDefinitionList []*ArgumentDefinition

func (l ArgumentDefinitionList) ByName(name string) *ArgumentDefinition {
	for _, arg := range l {
		if arg.Name == name {
			return arg
		}
	}
	return nil
}
