package schema

import "github.com/vektah/gqlparser/ast"

type ArgumentDefinition struct {
	*ast.ArgumentDefinition
	fieldDefinition *FieldDefinition
}

func (a *ArgumentDefinition) Type() *Type {
	return &Type{a.ArgumentDefinition.Type, a.fieldDefinition.definition.schema}
}

type ArgumentDefinitionList []*ArgumentDefinition
