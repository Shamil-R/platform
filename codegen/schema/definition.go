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
	schema     *Schema
	fields     FieldList
	directives DirectiveList
	mutations  ActionList
	queries    ActionList
	relations  ActionList
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
	if d.fields != nil {
		return d.fields
	}
	d.fields = make(FieldList, 0, len(d.Definition.Fields))
	for _, field := range d.Definition.Fields {
		if defaultFields[field.Name] == true {
			continue
		}
		d.fields = append(d.fields, &FieldDefinition{FieldDefinition: field, parent: d})
	}

	return d.fields
}

func (d *Definition) Directives() DirectiveList {
	if d.directives == nil {
		d.directives = make(DirectiveList, 0, len(d.Definition.Directives))
		for _, dir := range d.Definition.Directives {
			directive := &Directive{Directive: dir}
			d.directives = append(d.directives, directive)
		}
	}
	return d.directives
}

func (d *Definition) Mutations() ActionList {
	if d.mutations != nil {
		return d.mutations
	}

	d.mutations = make(ActionList, 0)
	mutation := d.schema.Mutation()
	if mutation == nil {
		return d.mutations
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
			d.mutations = append(d.mutations, action)
		}
	}

	return d.mutations
}

func (d *Definition) Queries() ActionList {
	if d.queries != nil {
		return d.queries
	}

	d.queries = make(ActionList, 0)
	query := d.schema.Query()
	if query == nil {
		return d.queries
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
			d.queries = append(d.queries, action)
		}
	}

	return d.queries
}

func (d *Definition) Relations() ActionList {
	if d.relations != nil {
		return d.relations
	}

	d.relations = make(ActionList, 0)
	for _, field := range d.Fields().Relations() {
		d.relations = append(d.relations, &Action{field, ACTION_RELATION})
	}

	return d.relations
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
