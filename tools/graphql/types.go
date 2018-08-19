package graphql

import (
	"strings"

	"github.com/jinzhu/inflection"
	"github.com/vektah/gqlparser/ast"
)

type Schema struct {
	*ast.Schema
	types []*Definition
}

func NewSchema(schema *ast.Schema) *Schema {
	return &Schema{
		types: toDefinitions(schema.Types),
	}
}

func (s *Schema) Types() []*Definition {
	return s.types
}

type Definition struct {
	*ast.Definition
	fields   []*Field
	Input    Input
	Mutation Mutation
	Query    Query
}

func (d *Definition) Fields() []*Field {
	return d.fields
}

type Field struct {
	*ast.FieldDefinition
}

func (f *Field) Describe() string {
	return f.Type.Dump()
}

type Input struct {
	Create      *CreateInput
	Update      *UpdateInput
	WhereUnique *WhereUniqueInput
	Where       *WhereInput
}

type CreateInput struct {
	*ast.Definition
}

func (i *CreateInput) Name() string {
	return i.Definition.Name + "CreateInput"
}

func (i *CreateInput) Fields() []*Field {
	fields := make([]*Field, 0, len(i.Definition.Fields))
	for _, f := range i.Definition.Fields {
		if f.Type.Name() != "ID" && f.Type.Elem == nil {
			fields = append(fields, &Field{f})
		}
	}
	return fields
}

type UpdateInput struct {
	*ast.Definition
}

func (i *UpdateInput) Name() string {
	return i.Definition.Name + "UpdateInput"
}

func (i *UpdateInput) Fields() []*Field {
	fields := make([]*Field, 0, len(i.Definition.Fields))
	for _, f := range i.Definition.Fields {
		if f.Type.Name() != "ID" && f.Type.Elem == nil {
			fields = append(fields, &Field{f})
		}
	}
	return fields
}

type WhereUniqueInput struct {
	*ast.Definition
}

func (i *WhereUniqueInput) Name() string {
	return i.Definition.Name + "WhereUniqueInput"
}

type WhereInput struct {
	*ast.Definition
}

func (i *WhereInput) Name() string {
	return i.Definition.Name + "WhereInput"
}

type Mutation struct {
	Create string
	Update string
	Delete string
}

type Query struct {
	Item string
	List string
}

func toDefinitions(list map[string]*ast.Definition) []*Definition {
	definitions := make([]*Definition, 0, len(list))
	for _, def := range list {
		if def.IsCompositeType() && !strings.HasPrefix(def.Name, "__") {
			definition := &Definition{
				Definition: def,
				fields:     toFields(def.Fields),
				Input: Input{
					Create:      &CreateInput{def},
					Update:      &UpdateInput{def},
					WhereUnique: &WhereUniqueInput{def},
					Where:       &WhereInput{def},
				},
				Mutation: Mutation{
					Create: "create" + def.Name,
					Update: "update" + def.Name,
					Delete: "delete" + def.Name,
				},
				Query: Query{
					Item: strings.ToLower(def.Name),
					List: inflection.Plural(strings.ToLower(def.Name)),
				},
			}
			definitions = append(definitions, definition)
		}
	}
	return definitions
}

func toFields(list ast.FieldList) []*Field {
	fields := make([]*Field, len(list))
	for i, def := range list {
		fields[i] = &Field{
			FieldDefinition: def,
		}
	}
	return fields
}
