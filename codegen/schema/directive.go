package schema

import "github.com/vektah/gqlparser/ast"

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

func (l DirectiveList) size() int {
	return len(l)
}

func (l DirectiveList) filter(filter directiveListFilter) DirectiveList {
	directives := make(DirectiveList, 0, len(l))
	for _, directive := range l {
		if filter(directive) {
			directives = append(directives, directive)
		}
	}
	return directives
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
