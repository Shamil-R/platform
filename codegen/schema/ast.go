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

func (s *Schema) Types() DefinitionList {
	definitions := make(DefinitionList, 0, len(s.Schema.Types))
	for _, def := range s.Schema.Types {
		if !strings.HasPrefix(def.Name, "__") {
			definitions = append(definitions, &Definition{def})
		}
	}
	return definitions
}

type Definition struct {
	*ast.Definition
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
		fields[i] = &FieldDefinition{field}
	}
	return fields
}

type DefinitionList []*Definition

type definitionListFilter func(def *Definition) bool

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
	if r.Size() == 0 {
		return nil
	}
	return r[0]
}

func (l DefinitionList) objects() DefinitionList {
	fn := func(def *Definition) bool {
		return def.IsObject()
	}
	return l.filter(fn)
}

func (l DefinitionList) Size() int {
	return len(l)
}

func (l DefinitionList) ForInput() DefinitionList {
	return l.objects()
}

func (l DefinitionList) ForMutation() DefinitionList {
	return l.objects()
}

func (l DefinitionList) ForQuery() DefinitionList {
	return l.objects()
}

func (l DefinitionList) Mutation() *Definition {
	fn := func(def *Definition) bool {
		return def.IsMutation()
	}
	return l.first(fn)
}

func (l DefinitionList) Query() *Definition {
	fn := func(def *Definition) bool {
		return def.IsQuery()
	}
	return l.first(fn)
}

type FieldDefinition struct {
	*ast.FieldDefinition
}

func (f *FieldDefinition) Arguments() ArgumentDefinitionList {
	arguments := make(ArgumentDefinitionList, len(f.FieldDefinition.Arguments))
	for i, arg := range f.FieldDefinition.Arguments {
		arguments[i] = &ArgumentDefinition{arg}
	}
	return arguments
}

func (f *FieldDefinition) Directives() DirectiveList {
	directives := make(DirectiveList, len(f.FieldDefinition.Directives))
	for i, directive := range f.FieldDefinition.Directives {
		directives[i] = &Directive{directive}
	}
	return directives
}

type FieldList []*FieldDefinition

type fieldListFilter func(field *FieldDefinition) bool

func (l FieldList) filter(filter fieldListFilter) FieldList {
	fields := make(FieldList, 0, len(l))
	for _, field := range l {
		if filter(field) {
			fields = append(fields, field)
		}
	}
	return fields
}

func (l FieldList) first(filter fieldListFilter) *FieldDefinition {
	r := l.filter(filter)
	if r.Size() == 0 {
		return nil
	}
	return r[0]
}

func (l FieldList) Size() int {
	return len(l)
}

func (l FieldList) ForObject() FieldList {
	return l
}

func (l FieldList) ForCreateInput() FieldList {
	fn := func(field *FieldDefinition) bool {
		return !field.Directives().HasIndentity()
	}
	return l.filter(fn)
}

func (l FieldList) ForUpdateInput() FieldList {
	return l.ForCreateInput()
}

func (l FieldList) ForWhereUniqueInput() FieldList {
	fn := func(field *FieldDefinition) bool {
		return field.Directives().HasIndentity()
	}
	return l.filter(fn)
}

func (l FieldList) ForWhereInput() FieldList {
	return l
}

func (l FieldList) ByName(name string) *FieldDefinition {
	fn := func(field *FieldDefinition) bool {
		return strings.ToLower(field.Name) == strings.ToLower(name)
	}
	return l.first(fn)
}

const (
	DIRECTIVE_PRIMARY   = "primary"
	DIRECTIVE_UNIQUE    = "unique"
	DIRECTIVE_INDENTITY = "identity"
	DIRECTIVE_VALIDATE  = "validate"
)

type ArgumentDefinition struct {
	*ast.ArgumentDefinition
}

func (a *ArgumentDefinition) IsSlice() bool {
	return a.Type.NamedType == "" && a.Type.Elem != nil
}

type ArgumentDefinitionList []*ArgumentDefinition

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
	directives := make(DirectiveList, 0, len(l))
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
	fn := func(directive *Directive) bool {
		return directive.IsPrimary()
	}
	return l.filter(fn).Size() > 0
}

func (l DirectiveList) HasUnique() bool {
	fn := func(directive *Directive) bool {
		return directive.IsUnique()
	}
	return l.filter(fn).Size() > 0
}

func (l DirectiveList) HasIndentity() bool {
	fn := func(directive *Directive) bool {
		return directive.IsIndentity()
	}
	return l.filter(fn).Size() > 0
}

func (l DirectiveList) HasValidate() bool {
	fn := func(directive *Directive) bool {
		return directive.IsValidate()
	}
	return l.filter(fn).Size() > 0
}

func (l DirectiveList) ForCreateInput() DirectiveList {
	fn := func(directive *Directive) bool {
		return directive.IsValidate()
	}
	return l.filter(fn)
}

func (l DirectiveList) ForUpdateInput() DirectiveList {
	return l.ForCreateInput()
}
