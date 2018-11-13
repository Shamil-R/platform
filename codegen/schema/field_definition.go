package schema

import (
	"github.com/vektah/gqlparser/ast"
)

type FieldDefinition struct {
	*ast.FieldDefinition
	parent   *Definition
	relation *FieldDefinition
	arguments ArgumentDefinitionList
	directives DirectiveList
}

func (f *FieldDefinition) Definition() *Definition {
	return f.parent
}

func (f *FieldDefinition) Type() *Type {
	return &Type{f.FieldDefinition.Type, f.parent.schema}
}

func (f *FieldDefinition) Arguments() ArgumentDefinitionList {
	if f.arguments != nil {
		return f.arguments
	}

	f.arguments = make(ArgumentDefinitionList, len(f.FieldDefinition.Arguments))
	for i, arg := range f.FieldDefinition.Arguments {
		f.arguments[i] = &ArgumentDefinition{arg, f}
	}

	return f.arguments
}

func (f *FieldDefinition) Directives() DirectiveList {
	if f.directives == nil {
		return f.directives
	}

	f.directives = make(DirectiveList, len(f.FieldDefinition.Directives))
	for i, directive := range f.FieldDefinition.Directives {
		f.directives[i] = &Directive{directive, nil}
	}

	return f.directives
}

func (f *FieldDefinition) Relation() *FieldDefinition {
	if f.relation == nil {
		if def := f.parent.schema.Types().ByType(f.Type()); def != nil {
			f.relation = def.Fields().ByType(f.parent.Name)
		}
	}
	return f.relation
}

type FieldList []*FieldDefinition

func (l FieldList) HasRelations() bool {
	return hasField(l, isRelation)
}

func (l FieldList) Primary() *FieldDefinition {
	return firstField(l, isPrimaryField)
}

func (l FieldList) RelationsOneToMany() FieldList {
	return filterFields(l, isOneToManyRelation)
}

func (l FieldList) RelationsManyToOne() FieldList {
	return filterFields(l, isManyToOneRelation)
}

func (l FieldList) Relations() FieldList {
	return filterFields(l, isRelation)
}

func (l FieldList) NotRelations() FieldList {
	return filterFields(l, notRelation)
}

func (l FieldList) ByName(name string) *FieldDefinition {
	fn := func(field *FieldDefinition) bool {
		return field.Name == name
	}
	return firstField(l, fn)
}

func (l FieldList) ByType(name string) *FieldDefinition {
	fn := func(field *FieldDefinition) bool {
		return field.Type().Name() == name
	}
	return firstField(l, fn)
}

func isPrimaryField(field *FieldDefinition) bool {
	return field.Directives().HasPrimary()
}

func isOneToManyRelation(field *FieldDefinition) bool {
	if !field.Type().IsSlice() {
		return false
	}
	if field.Relation() == nil {
		return false
	}
	return !field.Relation().Type().IsSlice()
}

func isManyToOneRelation(field *FieldDefinition) bool {
	if field.Type().IsSlice() {
		return false
	}
	if field.Relation() == nil {
		return false
	}
	return field.Relation().Type().IsSlice()
}

func isRelation(field *FieldDefinition) bool {
	return isOneToManyRelation(field) || isManyToOneRelation(field)
}

func notRelation(field *FieldDefinition) bool {
	return !isRelation(field)
}

type fieldFilter func(field *FieldDefinition) bool

func hasField(list FieldList, filter fieldFilter) bool {
	return firstField(list, filter) != nil
}

func firstField(list FieldList, filter fieldFilter) *FieldDefinition {
	for _, field := range list {
		if filter(field) {
			return field
		}
	}
	return nil
}

func filterFields(list FieldList, filter fieldFilter) FieldList {
	fields := make(FieldList, 0, len(list))
	for _, field := range list {
		if filter(field) {
			fields = append(fields, field)
		}
	}
	return fields
}
