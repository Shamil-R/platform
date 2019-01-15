package schema

import "github.com/vektah/gqlparser/ast"

type ArgumentDefinition struct {
	*ast.ArgumentDefinition
	fieldDefinition *FieldDefinition
	typeCache       *Type
}

func (a *ArgumentDefinition) Type() *Type {
	if a.typeCache == nil {
		a.typeCache = &Type{
			Type:  a.ArgumentDefinition.Type,
			types: a.fieldDefinition.parent.schema.Types(),
		}
	}
	return a.typeCache
}

type ArgumentDefinitionList []*ArgumentDefinition

type argumentDefinitionListFilter func(arg *ArgumentDefinition) bool

func (l ArgumentDefinitionList) size() int {
	return len(l)
}

func (l ArgumentDefinitionList) filter(filter argumentDefinitionListFilter) ArgumentDefinitionList {
	args := make(ArgumentDefinitionList, 0, len(l))
	for _, arg := range l {
		if filter(arg) {
			args = append(args, arg)
		}
	}
	return args
}

func (l ArgumentDefinitionList) first(filter argumentDefinitionListFilter) *ArgumentDefinition {
	r := l.filter(filter)
	if r.size() == 0 {
		return nil
	}
	return r[0]
}

func (l ArgumentDefinitionList) ByName(name string) *ArgumentDefinition {
	fn := func(arg *ArgumentDefinition) bool {
		return arg.Name == name
	}
	return l.first(fn)
}
