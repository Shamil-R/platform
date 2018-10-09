package schema

import "github.com/vektah/gqlparser/ast"

type FieldDefinition struct {
	*ast.FieldDefinition
	schema *Schema
}

func (f *FieldDefinition) Directives() DirectiveList {
	directives := make(DirectiveList, len(f.FieldDefinition.Directives))
	for i, directive := range f.FieldDefinition.Directives {
		directives[i] = &Directive{directive}
	}
	return directives
}

func (f *FieldDefinition) Arguments() ArgumentDefinitionList {
	arguments := make(ArgumentDefinitionList, len(f.FieldDefinition.Arguments))
	for i, arg := range f.FieldDefinition.Arguments {
		arguments[i] = &ArgumentDefinition{arg, f.schema}
	}
	return arguments
}

func (f *FieldDefinition) Type() *Type {
	return &Type{f.FieldDefinition.Type, f.schema}
}

type FieldList []*FieldDefinition

type fieldListFilter func(field *FieldDefinition) bool

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

func (l FieldList) size() int {
	return len(l)
}

func (l FieldList) Objects() FieldList {
	fn := func(field *FieldDefinition) bool {
		if field.Type().IsSlice() {
			return field.Type().Elem().IsObject()
		}
		return field.Type().IsObject()
	}
	return l.filter(fn)
}

func (l FieldList) ForObject() FieldList {
	return l
}

func (l FieldList) ForCreateInput() FieldList {
	fn := func(field *FieldDefinition) bool {
		return !field.Type().IsSlice() &&
			!field.Type().IsObject() &&
			!field.Directives().HasIndentity()
	}
	return l.filter(fn)
}

func (l FieldList) ForUpdateInput() FieldList {
	return l.ForCreateInput()
}

func (l FieldList) ForWhereUniqueInput() FieldList {
	fn := func(field *FieldDefinition) bool {
		return field.Directives().HasIndentity()
	}
	return l.filter(fn)
}

func (l FieldList) ForWhereInput() FieldList {
	return l.ForWhereUniqueInput()
}

func (l FieldList) ByName(name string) *FieldDefinition {
	fn := func(field *FieldDefinition) bool {
		return field.Name == name
	}
	return l.first(fn)
}
