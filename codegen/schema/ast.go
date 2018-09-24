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

func (f *FieldDefinition) Directives() DirectiveList {
	directives := make([]*Directive, len(f.FieldDefinition.Directives))
	for i, directive := range f.FieldDefinition.Directives {
		directives[i] = &Directive{directive}
	}
	return directives
}

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
		return !field.Directives().HasIndentity()
	}
	return l.filter(filter)
}

func (l FieldList) ForUpdateInput() FieldList {
	return l.ForCreateInput()
}

func (l FieldList) ForWhereUniqueInput() FieldList {
	filter := func(field *FieldDefinition) bool {
		return field.Directives().HasIndentity()
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
	DIRECTIVE_VALIDATE  = "validate"
)

type Directive struct {
	*ast.Directive
}

func (d *Directive) IsPrimary() bool {
	return d.Name == DIRECTIVE_PRIMARY
}

func (d *Directive) IsUnique() bool {
	return d.Name == DIRECTIVE_UNIQUE
}

func (d *Directive) IsIndentity() bool {
	return d.Name == DIRECTIVE_INDENTITY
}

func (d *Directive) IsValidate() bool {
	return d.Name == DIRECTIVE_VALIDATE
}

func (d *Directive) HasArguments() bool {
	return len(d.Arguments) > 0
}

type DirectiveList []*Directive

type directiveListFilter func(field *Directive) bool

func (l DirectiveList) filter(filter directiveListFilter) DirectiveList {
	directives := make([]*Directive, 0, len(l))
	for _, directive := range l {
		if filter(directive) {
			directives = append(directives, directive)
		}
	}
	return directives
}

func (l DirectiveList) Size() int {
	return len(l)
}

func (l DirectiveList) HasPrimary() bool {
	return l.filter(func(directive *Directive) bool {
		return directive.IsPrimary()
	}).Size() > 0
}

func (l DirectiveList) HasUnique() bool {
	return l.filter(func(directive *Directive) bool {
		return directive.IsUnique()
	}).Size() > 0
}

func (l DirectiveList) HasIndentity() bool {
	return l.filter(func(directive *Directive) bool {
		return directive.IsIndentity()
	}).Size() > 0
}

func (l DirectiveList) HasValidate() bool {
	return l.filter(func(directive *Directive) bool {
		return directive.IsValidate()
	}).Size() > 0
}

func (l DirectiveList) ForCreateInput() DirectiveList {
	return l.filter(func(directive *Directive) bool {
		return directive.IsValidate()
	})
}

func (l DirectiveList) ForUpdateInput() DirectiveList {
	return l.ForCreateInput()
}
