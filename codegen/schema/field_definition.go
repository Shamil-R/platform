package schema

import "github.com/vektah/gqlparser/ast"

type FieldDefinition struct {
	*ast.FieldDefinition
	definition *Definition
}

func (f *FieldDefinition) Definition() *Definition {
	return f.definition
}

func (f *FieldDefinition) Type() *Type {
	return &Type{f.FieldDefinition.Type, f.definition.schema}
}

func (f *FieldDefinition) Arguments() ArgumentDefinitionList {
	arguments := make(ArgumentDefinitionList, len(f.FieldDefinition.Arguments))
	for i, arg := range f.FieldDefinition.Arguments {
		arguments[i] = &ArgumentDefinition{arg, f}
	}
	return arguments
}

func (f *FieldDefinition) Directives() DirectiveList {
	directives := make(DirectiveList, len(f.FieldDefinition.Directives))
	for i, directive := range f.FieldDefinition.Directives {
		directives[i] = &Directive{directive}
	}
	return directives
}

type FieldList []*FieldDefinition

type fieldListFilter func(field *FieldDefinition) bool

func (l FieldList) size() int {
	return len(l)
}

func (l FieldList) filter(filter fieldListFilter) FieldList {
	fields := make(FieldList, 0, len(l))
	for _, field := range l {
		if filter(field) {
			fields = append(fields, field)
		}
	}
	return fields
}

func (l FieldList) first(filter fieldListFilter) *FieldDefinition {
	r := l.filter(filter)
	if r.size() == 0 {
		return nil
	}
	return r[0]
}

func (l FieldList) HasRelation() bool {
	fn := func(field *FieldDefinition) bool {
		return field.Type().IsRelation()
	}
	return l.filter(fn).size() > 0
}

func (l FieldList) Relations() FieldList {
	fn := func(field *FieldDefinition) bool {
		return field.Type().IsRelation()
	}
	return l.filter(fn)
}
