package schema

import "github.com/vektah/gqlparser/ast"

const (
	DIRECTIVE_PRIMARY   = "primary"
	DIRECTIVE_UNIQUE    = "unique"
	DIRECTIVE_INDENTITY = "identity"
	DIRECTIVE_VALIDATE  = "validate"
	DIRECTIVE_TABLE     = "table"
	DIRECTIVE_FIELD     = "field"
	DIRECTIVE_RELATION  = "relation"
)

type Directive struct {
	*ast.Directive
}

func (d *Directive) IsPrimary() bool {
	return isPrimaryDirective(d)
}

func (d *Directive) IsUnique() bool {
	return isUniqueDirective(d)
}

func (d *Directive) IsIndentity() bool {
	return isIndentityDirective(d)
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
	args := make(ArgumentList, len(f.Directive.Arguments))
	for i, arg := range f.Directive.Arguments {
		args[i] = &Argument{Argument: arg}
	}
	return args
}

type TableDirective struct {
	*Directive
	name *string
}

func (d *TableDirective) Name() string {
	if d.name == nil {
		d.name = &d.Arguments().ByName("name").Value().Raw
	}
	return *d.name
}

type FieldDirective struct {
	*Directive
	name *string
}

func (d *FieldDirective) Name() string {
	if d.name == nil {
		d.name = &d.Arguments().ByName("name").Value().Raw
	}
	return *d.name
}

type DirectiveList []*Directive

func (l DirectiveList) HasPrimary() bool {
	return hasDirective(l, isPrimaryDirective)
}

func (l DirectiveList) HasUnique() bool {
	return hasDirective(l, isUniqueDirective)
}

func (l DirectiveList) HasIndentity() bool {
	return hasDirective(l, isIndentityDirective)
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

func (l DirectiveList) ByName(name string) *Directive {
	return firstDirective(l, byNameDirective(name))
}

type directiveFilter func(directive *Directive) bool

func isPrimaryDirective(directive *Directive) bool {
	return directive.Name == DIRECTIVE_PRIMARY
}

func isUniqueDirective(directive *Directive) bool {
	return directive.Name == DIRECTIVE_UNIQUE
}

func isIndentityDirective(directive *Directive) bool {
	return directive.Name == DIRECTIVE_INDENTITY
}

func isValidateDirective(directive *Directive) bool {
	return directive.Name == DIRECTIVE_VALIDATE
}

func isTableDirective(directive *Directive) bool {
	return directive.Name == DIRECTIVE_TABLE
}

func isFieldDirective(directive *Directive) bool {
	return directive.Name == DIRECTIVE_FIELD
}

func isRelationDirective(directive *Directive) bool {
	return directive.Name == DIRECTIVE_RELATION
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
