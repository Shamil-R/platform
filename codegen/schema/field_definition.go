package schema

import (
	"github.com/vektah/gqlparser/ast"
)

type FieldDefinition struct {
	*ast.FieldDefinition
	parent     *Definition
	relation   *FieldDefinition
	typeCache  *Type
	arguments  ArgumentDefinitionList
	directives DirectiveList
}

func (f *FieldDefinition) Definition() *Definition {
	return f.parent
}

func (f *FieldDefinition) Type() *Type {
	if f.typeCache == nil {
		f.typeCache = &Type{
			Type:  f.FieldDefinition.Type,
			types: f.parent.schema.Types(),
		}
	}
	return f.typeCache
}

func (f *FieldDefinition) Arguments() ArgumentDefinitionList {
	if f.arguments == nil {
		l := len(f.FieldDefinition.Arguments)
		f.arguments = make(ArgumentDefinitionList, 0, l)
		for _, argument := range f.FieldDefinition.Arguments {
			arg := &ArgumentDefinition{
				ArgumentDefinition: argument,
				fieldDefinition:    f,
			}
			f.arguments = append(f.arguments, arg)
		}
	}
	return f.arguments
}

func (f *FieldDefinition) Directives() DirectiveList {
	if f.directives == nil {
		f.directives = make(DirectiveList, 0, len(f.FieldDefinition.Directives))
		for _, directive := range f.FieldDefinition.Directives {
			dir := &Directive{Directive: directive}
			f.directives = append(f.directives, dir)
		}
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

// type Relation struct {
// 	schema               *Schema
// 	owner                *FieldDefinition
// 	definitionCache      *Definition
// 	fieldDefinitionCache *FieldDefinition
// }

// func (r *Relation) Definition() *Definition {
// 	if r.definitionCache == nil {
// 		r.definitionCache = r.schema.Types().ByType(r.owner.Type())
// 	}
// 	return r.definitionCache
// }

// func (r *Relation) FieldDefinition() *FieldDefinition {
// 	return r.fieldDefinitionCache
// }

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
