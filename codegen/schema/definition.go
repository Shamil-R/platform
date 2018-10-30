package schema

import (
	"github.com/huandu/xstrings"
	"github.com/jinzhu/inflection"
	"github.com/vektah/gqlparser/ast"
)

var defaultFields map[string]bool = map[string]bool{
	"createdAt": true,
	"updatedAt": true,
	"deletedAt": true,
}

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
	fields := make(FieldList, 0, len(d.Definition.Fields))
	for _, field := range d.Definition.Fields {
		if defaultFields[field.Name] == true {
			continue
		}
		fields = append(fields, &FieldDefinition{FieldDefinition: field, parent: d})
	}
	return fields
}

func (d *Definition) Directives() DirectiveList {
	directives := make(DirectiveList, len(d.Definition.Directives))

	for i, d := range d.Definition.Directives {
		directives[i] = &Directive{d}
	}
	return directives
}

func (d *Definition) Mutations() ActionList {
	actions := make(ActionList, 0)
	mutation := d.schema.Mutation()
	if mutation == nil {
		return actions
	}
	checks := map[string]string{
		ACTION_CREATE + d.Name: ACTION_CREATE,
		ACTION_UPDATE + d.Name: ACTION_UPDATE,
		ACTION_DELETE + d.Name: ACTION_DELETE,
		ACTION_UPSERT + d.Name: ACTION_UPSERT,
	}
	for _, field := range mutation.Fields() {
		if act, ok := checks[field.Name]; ok {
			action := &Action{&FieldDefinition{FieldDefinition: field.FieldDefinition, parent: d}, act}
			actions = append(actions, action)
		}
	}
	return actions
}

func (d *Definition) Queries() ActionList {
	actions := make(ActionList, 0)
	query := d.schema.Query()
	if query == nil {
		return actions
	}
	item := xstrings.FirstRuneToLower(d.Name)
	collection := inflection.Plural(item)
	checks := map[string]string{
		item:       ACTION_ITEM,
		collection: ACTION_COLLECTION,
	}
	for _, field := range query.Fields() {
		if act, ok := checks[field.Name]; ok {
			action := &Action{&FieldDefinition{FieldDefinition: field.FieldDefinition, parent: d}, act}
			actions = append(actions, action)
		}
	}
	return actions
}

func (d *Definition) Relations() ActionList {
	acts := make(ActionList, 0)
	for _, field := range d.Fields().Relations() {
		acts = append(acts, &Action{field, ACTION_RELATION})
	}
	return acts
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

func (l DefinitionList) ByName(name string) *Definition {
	fn := func(def *Definition) bool {
		return def.Name == name
	}
	return l.first(fn)
}

func (l DefinitionList) ByType(t *Type) *Definition {
	fn := func(def *Definition) bool {
		return def.Name == t.Name()
	}
	return l.first(fn)
}
