package schema

import (
	"github.com/vektah/gqlparser/ast"
)

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
	args := make(ArgumentDefinitionList, len(f.FieldDefinition.Arguments))
	for i, arg := range f.FieldDefinition.Arguments {
		args[i] = &ArgumentDefinition{arg, f}
	}
	return args
}

func (f *FieldDefinition) Directives() DirectiveList {
	directives := make(DirectiveList, len(f.FieldDefinition.Directives))
	for i, directive := range f.FieldDefinition.Directives {
		directives[i] = &Directive{directive}
	}
	return directives
}

func (f *FieldDefinition) Relation() *Relation {
	return &Relation{f}
}

func (f *FieldDefinition) IsRelation() bool {
	return f.Relation().IsRelation()
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

func (l FieldList) HasRelations() bool {
	fn := func(field *FieldDefinition) bool {
		return field.IsRelation()
	}
	return l.filter(fn).size() > 0
}

func (l FieldList) Relations() RelationList {
	fn := func(field *FieldDefinition) bool {
		return field.IsRelation()
	}
	return RelationList(l.filter(fn))
}

func (l FieldList) NotRelations() FieldList {
	fn := func(field *FieldDefinition) bool {
		return !field.IsRelation()
	}
	return l.filter(fn)
}

func (l FieldList) ByNameType(name string) *FieldDefinition {
	fn := func(field *FieldDefinition) bool {
		return field.Type().Name() == name
	}
	return l.first(fn)
}

func hasRelation(field *FieldDefinition) bool {
	t := field.Type()
	return t.schema.Types().Objects().Contains(t)
}

type Relation struct {
	*FieldDefinition
}

func (r *Relation) Field() *FieldDefinition {
	def := r.definition.schema.Types().ByType(r.Type())
	if def == nil {
		return nil
	}
	return def.Fields().ByNameType(r.definition.Name)
}

func (r *Relation) IsOneToMany() bool {
	if !r.Type().IsSlice() {
		return false
	}
	field := r.Field()
	if field == nil {
		return false
	}
	return !field.Type().IsSlice()
}

func (r *Relation) IsManyToOne() bool {
	if r.Type().IsSlice() {
		return false
	}
	field := r.Field()
	if field == nil {
		return false
	}
	return field.Type().IsSlice()
}

func (r *Relation) IsRelation() bool {
	return r.IsOneToMany() || r.IsManyToOne()
}

type RelationList FieldList

func (l RelationList) OneToMany() FieldList {
	fn := func(field *FieldDefinition) bool {
		return field.Relation().IsOneToMany()
	}
	return FieldList(l).filter(fn)
}

func (l RelationList) ManyToOne() FieldList {
	fn := func(field *FieldDefinition) bool {
		return field.Relation().IsManyToOne()
	}
	return FieldList(l).filter(fn)
}
