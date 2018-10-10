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
		fields[i] = &FieldDefinition{field, d.schema, d}
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
		return d.schema.Types().Objects().ByName(field.Type().Name()) != nil
	}
	return d.Fields().filter(fn)
}

/* func (d *Definition) Mutations() ActionList {
	mutation := d.schema.Mutation()
	if mutation == nil {
		return ActionList{}
	}
	checks := map[string]string{
		ACTION_CREATE + d.Name: ACTION_CREATE,
		ACTION_UPDATE + d.Name: ACTION_UPDATE,
		ACTION_DELETE + d.Name: ACTION_DELETE,
	}
	actions := make(ActionList, 0, len(mutation.Fields()))
	for _, field := range mutation.Fields() {
		if action, ok := checks[field.Name]; ok {
			actions = append(actions, &Action{action, field, d})
		}

	}
	return actions
}

func (d *Definition) Queries() ActionList {
	query := d.schema.Query()
	if query == nil {
		return ActionList{}
	}
	item := xstrings.FirstRuneToLower(d.Name)
	collection := inflection.Plural(item)
	checks := map[string]string{
		item:       ACTION_ITEM,
		collection: ACTION_COLLECTION,
	}
	actions := make(ActionList, 0, len(query.Fields()))
	for _, field := range query.Fields() {
		if action, ok := checks[field.Name]; ok {
			actions = append(actions, &Action{action, field, d})
		}
	}
	return actions
}

func (d *Definition) Relations() ActionList {
	actions := make(ActionList, 0, len(d.Fields()))
	for _, field := range d.Fields() {
		if field.Type().IsSlice() {
			actions = append(actions, &Action{ACTION_RELATION, field, d})
		}
	}
	return actions
}

func (d *Definition) Actions() ActionList {
	actions := append(d.Mutations(), d.Queries()...)
	actions = append(actions, d.Relations()...)
	return actions
} */

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

func (l DefinitionList) Objects() DefinitionList {
	fn := func(def *Definition) bool {
		return def.IsObject()
	}
	return l.filter(fn)
}

// func (l DefinitionList) ForObject() DefinitionList {
// 	return l.Objects()
// }

func (l DefinitionList) Enums() DefinitionList {
	fn := func(def *Definition) bool {
		return def.IsEnum()
	}
	return l.filter(fn)
}

// func (l DefinitionList) ForInput() DefinitionList {
// 	return l.Objects()
// }

// func (l DefinitionList) ForMutation() DefinitionList {
// 	return l.Objects()
// }

// func (l DefinitionList) ForQuery() DefinitionList {
// 	return l.Objects()
// }

// func (l DefinitionList) ForAction() DefinitionList {
// 	return l.Objects()
// }

// func (l DefinitionList) WithRelations() DefinitionList {
// 	fn := func(def *Definition) bool {
// 		return len(def.Fields().Objects()) > 0
// 	}
// 	return l.filter(fn)
// }

func (l DefinitionList) ByName(name string) *Definition {
	fn := func(def *Definition) bool {
		return def.Name == name
	}
	return l.first(fn)
}

func (l DefinitionList) ContainsType(t *Type) bool {
	fn := func(def *Definition) bool {
		return def.Name == t.Name()
	}
	return l.filter(fn).size() > 0
}
