package schema

import "github.com/vektah/gqlparser/ast"

type Field struct {
	*ast.Field
	argumentsCache        ArgumentList
	directivesCache       DirectiveList
	selectionSetCache     *SelectionSet
	definitionCache       *FieldDefinition
	objectDefinitionCache *Definition
}

func (f *Field) Arguments() ArgumentList {
	if f.argumentsCache == nil {
		f.argumentsCache = make(ArgumentList, 0, len(f.Field.Arguments))
		for _, arg := range f.Field.Arguments {
			f.argumentsCache = append(f.argumentsCache, &Argument{Argument: arg})
		}
	}
	return f.argumentsCache
}

func (f *Field) Directives() DirectiveList {
	if f.directivesCache == nil {
		f.directivesCache = make(DirectiveList, 0, len(f.Field.Directives))
		for _, directive := range f.Field.Directives {
			dir := &Directive{Directive: directive}
			f.directivesCache = append(f.directivesCache, dir)
		}
	}
	return f.directivesCache
}

func (f *Field) SelectionSet() *SelectionSet {
	if f.selectionSetCache == nil {
		f.selectionSetCache = &SelectionSet{SelectionSet: f.Field.SelectionSet}
	}
	return f.selectionSetCache
}

func (f *Field) Definition() *FieldDefinition {
	if f.definitionCache == nil {
		f.definitionCache = &FieldDefinition{FieldDefinition: f.Field.Definition}
	}
	return f.definitionCache
}

func (f *Field) ObjectDefinition() *Definition {
	if f.objectDefinitionCache == nil {
		f.objectDefinitionCache = &Definition{Definition: f.Field.ObjectDefinition}
	}
	return f.objectDefinitionCache
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
