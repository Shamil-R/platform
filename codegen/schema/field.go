package schema

import "github.com/vektah/gqlparser/ast"

type Field struct {
	*ast.Field
	selectionSet     *SelectionSet
	definition       *FieldDefinition
	objectDefinition *Definition
}

func (f *Field) SelectionSet() *SelectionSet {
	if f.selectionSet == nil {
		f.selectionSet = &SelectionSet{SelectionSet: f.Field.SelectionSet}
	}
	return f.selectionSet
}

func (f *Field) Definition() *FieldDefinition {
	if f.definition == nil {
		f.definition = &FieldDefinition{FieldDefinition: f.Field.Definition}
	}
	return f.definition
}

func (f *Field) ObjectDefinition() *Definition {
	if f.objectDefinition == nil {
		f.objectDefinition = &Definition{Definition: f.Field.ObjectDefinition}
	}
	return f.objectDefinition
}

type SelectionSet struct {
	ast.SelectionSet
	fields []*Field
}

func (s SelectionSet) Fields() []*Field {
	if s.fields == nil {
		s.fields = make([]*Field, 0, len(s.SelectionSet))
		for _, sel := range s.SelectionSet {
			if f, ok := sel.(*ast.Field); ok {
				field := &Field{Field: f}
				s.fields = append(s.fields, field)
			}
		}
	}
	return s.fields
}
