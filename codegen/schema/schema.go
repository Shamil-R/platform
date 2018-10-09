package schema

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/gobuffalo/packr"
	"github.com/huandu/xstrings"
	"github.com/jinzhu/inflection"
	"github.com/vektah/gqlparser"
	"github.com/vektah/gqlparser/ast"
)

type Schema struct {
	*ast.Schema
}

func (s *Schema) Types() DefinitionList {
	definitions := make(DefinitionList, 0, len(s.Schema.Types))
	for _, def := range s.Schema.Types {
		if !strings.HasPrefix(def.Name, "__") {
			definitions = append(definitions, &Definition{def, s})
		}
	}
	return definitions
}

func (s *Schema) Mutation() *Definition {
	if s.Schema.Mutation == nil {
		return nil
	}
	return &Definition{s.Schema.Mutation, s}
}

func (s *Schema) Query() *Definition {
	if s.Schema.Query == nil {
		return nil
	}
	return &Definition{s.Schema.Query, s}
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
	fields := make(FieldList, len(d.Definition.Fields))
	for i, field := range d.Definition.Fields {
		fields[i] = &FieldDefinition{field, d.schema}
	}
	return fields
}

func (d *Definition) Mutations() ActionList {
	mutation := d.schema.Mutation()
	if mutation == nil {
		return ActionList{}
	}
	checks := map[string]string{
		fmt.Sprintf("create%s", d.Name): "create",
		fmt.Sprintf("update%s", d.Name): "update",
		fmt.Sprintf("delete%s", d.Name): "delete",
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
		item:       "item",
		collection: "collection",
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
			actions = append(actions, &Action{"relation", field, d})
		}
	}
	return actions
}

func (d *Definition) Actions() ActionList {
	return append(d.Mutations(), d.Queries()...)
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
	if r.size() == 0 {
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

func (l DefinitionList) size() int {
	return len(l)
}

func (l DefinitionList) ForObject() DefinitionList {
	return l.objects()
}

func (l DefinitionList) ForEnum() DefinitionList {
	fn := func(def *Definition) bool {
		return def.IsEnum()
	}
	return l.filter(fn)
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

func (l DefinitionList) ForAction() DefinitionList {
	return l.objects()
}

func (l DefinitionList) ByName(name string) *Definition {
	fn := func(def *Definition) bool {
		return def.Name == name
	}
	return l.first(fn)
}

type FieldDefinition struct {
	*ast.FieldDefinition
	schema *Schema
}

func (f *FieldDefinition) Directives() DirectiveList {
	directives := make(DirectiveList, len(f.FieldDefinition.Directives))
	for i, directive := range f.FieldDefinition.Directives {
		directives[i] = &Directive{directive}
	}
	return directives
}

func (f *FieldDefinition) Arguments() ArgumentDefinitionList {
	arguments := make(ArgumentDefinitionList, len(f.FieldDefinition.Arguments))
	for i, arg := range f.FieldDefinition.Arguments {
		arguments[i] = &ArgumentDefinition{arg, f.schema}
	}
	return arguments
}

func (f *FieldDefinition) Type() *Type {
	return &Type{f.FieldDefinition.Type, f.schema}
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
	if r.size() == 0 {
		return nil
	}
	return r[0]
}

func (l FieldList) size() int {
	return len(l)
}

func (l FieldList) ForObject() FieldList {
	return l
}

func (l FieldList) ForCreateInput() FieldList {
	fn := func(field *FieldDefinition) bool {
		return !field.Type().IsSlice() &&
			!field.Type().IsObject() &&
			!field.Directives().HasIndentity()
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
	return l.ForWhereUniqueInput()
}

func (l FieldList) ByName(name string) *FieldDefinition {
	fn := func(field *FieldDefinition) bool {
		return field.Name == name
	}
	return l.first(fn)
}

type Type struct {
	*ast.Type
	schema *Schema
}

func (t *Type) IsObject() bool {
	return t.schema.Types().ForObject().ByName(t.NamedType) != nil
}

func (t *Type) IsSlice() bool {
	return t.NamedType == "" && t.Elem() != nil
}

func (t *Type) Elem() *Type {
	if t.Type.Elem == nil {
		return nil
	}
	return &Type{t.Type.Elem, t.schema}
}

type ArgumentDefinition struct {
	*ast.ArgumentDefinition
	schema *Schema
}

func (a *ArgumentDefinition) Type() *Type {
	return &Type{a.ArgumentDefinition.Type, a.schema}
}

type ArgumentDefinitionList []*ArgumentDefinition

func (l ArgumentDefinitionList) ByName(name string) *ArgumentDefinition {
	for _, arg := range l {
		if arg.Name == name {
			return arg
		}
	}
	return nil
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
	directives := make(DirectiveList, 0, len(l))
	for _, directive := range l {
		if filter(directive) {
			directives = append(directives, directive)
		}
	}
	return directives
}

func (l DirectiveList) size() int {
	return len(l)
}

func (l DirectiveList) HasPrimary() bool {
	fn := func(directive *Directive) bool {
		return directive.IsPrimary()
	}
	return l.filter(fn).size() > 0
}

func (l DirectiveList) HasUnique() bool {
	fn := func(directive *Directive) bool {
		return directive.IsUnique()
	}
	return l.filter(fn).size() > 0
}

func (l DirectiveList) HasIndentity() bool {
	fn := func(directive *Directive) bool {
		return directive.IsIndentity()
	}
	return l.filter(fn).size() > 0
}

func (l DirectiveList) HasValidate() bool {
	fn := func(directive *Directive) bool {
		return directive.IsValidate()
	}
	return l.filter(fn).size() > 0
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

type Action struct {
	Action          string
	FieldDefinition *FieldDefinition
	Definition      *Definition
}

type ActionList []*Action

type actionListFilter func(field *Action) bool

func (l ActionList) size() int {
	return len(l)
}

func (l ActionList) filter(filter actionListFilter) ActionList {
	actions := make(ActionList, 0, len(l))
	for _, action := range l {
		if filter(action) {
			actions = append(actions, action)
		}
	}
	return actions
}

func (l ActionList) first(filter actionListFilter) *Action {
	r := l.filter(filter)
	if r.size() == 0 {
		return nil
	}
	return r[0]
}

func (l ActionList) ByAction(action string) *Action {
	fn := func(a *Action) bool {
		return a.Action == action
	}
	return l.first(fn)
}

func LoadSchemaRaw(path string) (string, error) {
	box := packr.NewBox("./graphql")

	directivesRaw := box.Bytes("directives.graphql")

	buf := bytes.NewBuffer(directivesRaw)

	schemaRaw, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	if _, err := buf.Write(schemaRaw); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func LoadSchema(path string) (*Schema, error) {
	schemaRaw, err := LoadSchemaRaw(path)
	if err != nil {
		return nil, err
	}

	source := &ast.Source{
		Name:  "schema",
		Input: schemaRaw,
	}

	schema, gqlerr := gqlparser.LoadSchema(source)
	if gqlerr != nil {
		return nil, gqlerr
	}

	return &Schema{schema}, nil
}
