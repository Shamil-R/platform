package schema

import (
	"github.com/huandu/xstrings"
	"github.com/jinzhu/inflection"
	"github.com/vektah/gqlparser/ast"
)

type Definition struct {
	*ast.Definition
	schema *Schema
}

func (d *Definition) IsMutation() bool {
	return d.Name == "Mutation"
}

func (d *Definition) IsQuery() bool {
	return d.Name == "Query"
}

func (d *Definition) IsObject() bool {
	return !d.IsMutation() && !d.IsQuery() && d.Kind == ast.Object
}

func (d *Definition) IsEnum() bool {
	return d.Kind == ast.Enum
}

func (d *Definition) Fields() FieldList {
	fields := make(FieldList, len(d.Definition.Fields))
	for i, field := range d.Definition.Fields {
		fields[i] = &FieldDefinition{field, d}
	}
	return fields
}

func filterFieldList(def *Definition, checks map[string]bool) FieldList {
	fields := make(FieldList, 0)
	if def != nil {
		for _, field := range def.Fields() {
			if _, ok := checks[field.Name]; ok {
				fields = append(fields, field)
			}
		}
	}
	return fields
}

func (d *Definition) Mutations() FieldList {
	checks := map[string]bool{
		"create" + d.Name: true,
		"update" + d.Name: true,
		"delete" + d.Name: true,
	}
	return filterFieldList(d.schema.Mutation(), checks)
}

func (d *Definition) Queries() FieldList {
	item := xstrings.FirstRuneToLower(d.Name)
	collection := inflection.Plural(item)
	checks := map[string]bool{
		item:       true,
		collection: true,
	}
	return filterFieldList(d.schema.Query(), checks)
}

func (d *Definition) Relations() FieldList {
	fn := func(field *FieldDefinition) bool {
		return d.schema.Types().Objects().Contains(field.Type())
	}
	return d.Fields().filter(fn)
}

type DefinitionList []*Definition

type definitionListFilter func(def *Definition) bool

func (l DefinitionList) size() int {
	return len(l)
}

func (l DefinitionList) filter(filter definitionListFilter) DefinitionList {
	definitions := make(DefinitionList, 0, len(l))
	for _, def := range l {
		if filter(def) {
			definitions = append(definitions, def)
		}
	}
	return definitions
}

func (l DefinitionList) first(filter definitionListFilter) *Definition {
	r := l.filter(filter)
	if r.size() == 0 {
		return nil
	}
	return r[0]
}

func (l DefinitionList) Contains(t *Type) bool {
	fn := func(def *Definition) bool {
		return def.Name == t.Name()
	}
	return l.filter(fn).size() > 0
}

func (l DefinitionList) Objects() DefinitionList {
	fn := func(def *Definition) bool {
		return def.IsObject()
	}
	return l.filter(fn)
}

func (l DefinitionList) Enums() DefinitionList {
	fn := func(def *Definition) bool {
		return def.IsEnum()
	}
	return l.filter(fn)
}
