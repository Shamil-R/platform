package schema

import "github.com/vektah/gqlparser/ast"

type Schema struct {
	*ast.Schema
	types []*Definition
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
