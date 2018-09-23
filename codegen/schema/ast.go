package schema

import (
	"strings"

	"github.com/vektah/gqlparser/ast"
)

type Schema struct {
	*ast.Schema
}

func NewSchema(schema *ast.Schema) *Schema {
	return &Schema{schema}
}

func (s *Schema) ObjectTypes() []*Definition {
	return s.definitions(ast.Object)
}

func (s *Schema) EnumTypes() []*Definition {
	return s.definitions(ast.Enum)
}

func (s *Schema) InputTypes() []*Definition {
	return s.ObjectTypes()
}

func (s *Schema) MutationTypes() []*Definition {
	return s.ObjectTypes()
}

func (s *Schema) QueryTypes() []*Definition {
	return s.ObjectTypes()
}

func (s *Schema) definitions(kind ast.DefinitionKind) []*Definition {
	definitions := make([]*Definition, 0)
	for _, def := range s.Schema.Types {
		if !strings.HasPrefix(def.Name, "__") && def.Kind == kind {
			definitions = append(definitions, &Definition{def})
		}
	}
	return definitions
}

type Definition struct {
	*ast.Definition
}

func (d *Definition) Fields() FieldList {
	fields := make([]*FieldDefinition, len(d.Definition.Fields))
	for i, field := range d.Definition.Fields {
		fields[i] = &FieldDefinition{field}
	}
	return fields
}

type FieldDefinition struct {
	*ast.FieldDefinition
}

func (f *FieldDefinition) hasDirective(name string) bool {
	return f.FieldDefinition.Directives.ForName(name) != nil
}

func (f *FieldDefinition) HasPrimaryDirective() bool {
	return f.hasDirective(DIRECTIVE_PRIMARY)
}

func (f *FieldDefinition) HasUniqueDirective() bool {
	return f.hasDirective(DIRECTIVE_UNIQUE)
}

func (f *FieldDefinition) HasIndentityDirective() bool {
	return f.hasDirective(DIRECTIVE_INDENTITY)
}

/* func (f *FieldDefinition) Directives() []*Directive {
	directives := make([]*Directive, len(f.FieldDefinition.Directives))
	for i, directive := range f.FieldDefinition.Directives {
		directives[i] = &Directive{directive}
	}
	return directives
} */

type FieldList []*FieldDefinition

type fieldListFilter func(field *FieldDefinition) bool

func (l FieldList) filter(filter fieldListFilter) FieldList {
	fields := make([]*FieldDefinition, 0, len(l))
	for _, field := range l {
		if filter(field) {
			fields = append(fields, field)
		}
	}
	return fields
}

func (l FieldList) ForObject() FieldList {
	return l
}

func (l FieldList) ForCreateInput() FieldList {
	filter := func(field *FieldDefinition) bool {
		return !field.HasIndentityDirective()
	}
	return l.filter(filter)
}

func (l FieldList) ForUpdateInput() FieldList {
	return l.ForCreateInput()
}

func (l FieldList) ForWhereUniqueInput() FieldList {
	filter := func(field *FieldDefinition) bool {
		return field.HasIndentityDirective()
	}
	return l.filter(filter)
}

func (l FieldList) ForWhereInput() FieldList {
	return l
}

const (
	DIRECTIVE_PRIMARY   = "primary"
	DIRECTIVE_UNIQUE    = "unique"
	DIRECTIVE_INDENTITY = "identity"
)

/* type Directive struct {
	*ast.Directive
} */
