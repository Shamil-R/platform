package schema

import (
	"github.com/vektah/gqlparser/ast"
)

type EnumValueDefinition struct {
	*ast.EnumValueDefinition
	directives DirectiveList
}

func (f *EnumValueDefinition) Directives() DirectiveList {
	if f.directives == nil {
		f.directives = make(DirectiveList, 0, len(f.EnumValueDefinition.Directives))
		for _, directive := range f.EnumValueDefinition.Directives {
			dir := &Directive{Directive: directive}
			f.directives = append(f.directives, dir)
		}
	}
	return f.directives
}

type EnumValueList []*EnumValueDefinition

type fieldEnumFilter func(field *EnumValueDefinition) bool

func (l EnumValueList) ByName(name string) *EnumValueDefinition {
	fn := func(field *EnumValueDefinition) bool {
		return field.Name == name
	}
	return firstEnumField(l, fn)
}

func firstEnumField(list EnumValueList, filter fieldEnumFilter) *EnumValueDefinition {
	for _, field := range list {
		if filter(field) {
			return field
		}
	}
	return nil
}

