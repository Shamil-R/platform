package schema

import "github.com/vektah/gqlparser/ast"

const (
	DirectivePrimary  = "primary"
	DirectiveUnique   = "unique"
	DirectiveIdentity = "identity"
	DirectiveValidate = "validate"
	DirectiveTable    = "table"
	DirectiveField    = "field"
	DirectiveRelation = "relation"
	DirectiveInput    = "input"

	InputDirectiveCreateOneWithout = "create_one_without"
)

type Directives interface {
	Directives() DirectiveList
}

type Directive struct {
	*ast.Directive
	arguments ArgumentList
}

func (d *Directive) IsPrimary() bool {
	return isPrimaryDirective(d)
}

func (d *Directive) IsUnique() bool {
	return isUniqueDirective(d)
}

func (d *Directive) IsIndentity() bool {
	return isIdentityDirective(d)
}

func (d *Directive) IsValidate() bool {
	return isValidateDirective(d)
}

func (d *Directive) IsTable() bool {
	return isTableDirective(d)
}

func (d *Directive) IsField() bool {
	return isFieldDirective(d)
}

func (f *Directive) Arguments() ArgumentList {
	if f.arguments != nil {
		return f.arguments
	}

	f.arguments = make(ArgumentList, len(f.Directive.Arguments))
	for i, arg := range f.Directive.Arguments {
		f.arguments[i] = &Argument{Argument: arg}
	}

	return f.arguments
}

type ValidateDirective struct {
	*Directive
}

type TableDirective struct {
	*Directive
	argName *string
}

func (d *TableDirective) ArgName() *string {
	if d.argName == nil {
		arg := d.Arguments().ByName("name")
		if arg == nil {
			return nil
		}
		val := arg.Value()
		if val == nil {
			return nil
		}
		d.argName = &val.Raw
	}
	return d.argName
}

type FieldDirective struct {
	*Directive
	argName *string
}

func (d *FieldDirective) ArgName() string {
	if d.argName == nil {
		d.argName = &d.Arguments().ByName("name").Value().Raw
	}
	return *d.argName
}

type RelationDirective struct {
	*Directive
	argObject     *string
	argField      *string
	argForeignKey *string
}

func (d *RelationDirective) ArgObject() string {
	if d.argObject == nil {
		d.argObject = &d.Arguments().ByName("object").Value().Raw
	}
	return *d.argObject
}

func (d *RelationDirective) ArgField() string {
	if d.argField == nil {
		d.argField = &d.Arguments().ByName("field").Value().Raw
	}
	return *d.argField
}

func (d *RelationDirective) ArgForeignKey() string {
	if d.argForeignKey == nil {
		d.argForeignKey = &d.Arguments().ByName("foreignKey").Value().Raw
	}
	return *d.argForeignKey
}

type InputDirective struct {
	*Directive
}

func (d *InputDirective) IsCreateOneWithout() bool {
	return d.Arguments().ByName("name").Value().Raw == InputDirectiveCreateOneWithout
}

type DirectiveList []*Directive

func (l DirectiveList) HasPrimary() bool {
	return hasDirective(l, isPrimaryDirective)
}

func (l DirectiveList) HasUnique() bool {
	return hasDirective(l, isUniqueDirective)
}

func (l DirectiveList) HasIndentity() bool {
	return hasDirective(l, isIdentityDirective)
}

func (l DirectiveList) HasValidate() bool {
	return hasDirective(l, isValidateDirective)
}

func (l DirectiveList) HasTable() bool {
	return hasDirective(l, isTableDirective)
}

func (l DirectiveList) HasField() bool {
	return hasDirective(l, isFieldDirective)
}

func (l DirectiveList) Primary() *Directive {
	return firstDirective(l, isPrimaryDirective)
}

func (l DirectiveList) Unique() *Directive {
	return firstDirective(l, isUniqueDirective)
}

func (l DirectiveList) Identity() *Directive {
	return firstDirective(l, isIdentityDirective)
}

func (l DirectiveList) Validate() *ValidateDirective {
	directive := firstDirective(l, isValidateDirective)
	if directive == nil {
		return nil
	}
	return &ValidateDirective{Directive: directive}
}

func (l DirectiveList) Table() *TableDirective {
	directive := firstDirective(l, isTableDirective)
	if directive == nil {
		return nil
	}
	return &TableDirective{Directive: directive}
}

func (l DirectiveList) Field() *FieldDirective {
	directive := firstDirective(l, isFieldDirective)
	if directive == nil {
		return nil
	}
	return &FieldDirective{Directive: directive}
}

func (l DirectiveList) Relation() *RelationDirective {
	directive := firstDirective(l, isRelationDirective)
	if directive == nil {
		return nil
	}
	return &RelationDirective{Directive: directive}
}

func (l DirectiveList) Input() *InputDirective {
	directive := firstDirective(l, isInputDirective)
	if directive == nil {
		return nil
	}
	return &InputDirective{Directive: directive}
}

func (l DirectiveList) ByName(name string) *Directive {
	return firstDirective(l, byNameDirective(name))
}

type directiveFilter func(directive *Directive) bool

func isPrimaryDirective(directive *Directive) bool {
	return directive.Name == DirectivePrimary
}

func isUniqueDirective(directive *Directive) bool {
	return directive.Name == DirectiveUnique
}

func isIdentityDirective(directive *Directive) bool {
	return directive.Name == DirectiveIdentity
}

func isValidateDirective(directive *Directive) bool {
	return directive.Name == DirectiveValidate
}

func isTableDirective(directive *Directive) bool {
	return directive.Name == DirectiveTable
}

func isFieldDirective(directive *Directive) bool {
	return directive.Name == DirectiveField
}

func isRelationDirective(directive *Directive) bool {
	return directive.Name == DirectiveRelation
}

func isInputDirective(directive *Directive) bool {
	return directive.Name == DirectiveInput
}

func byNameDirective(name string) directiveFilter {
	return func(directive *Directive) bool {
		return directive.Name == name
	}
}

func hasDirective(list DirectiveList, filter directiveFilter) bool {
	for _, directive := range list {
		if filter(directive) {
			return true
		}
	}
	return false
}

func firstDirective(list DirectiveList, filter directiveFilter) *Directive {
	for _, directive := range list {
		if filter(directive) {
			return directive
		}
	}
	return nil
}
