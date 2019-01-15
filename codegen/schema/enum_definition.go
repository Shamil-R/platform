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

type valueEnumFilter func(value *EnumValueDefinition) bool

func (l EnumValueList) ByName(name string) *EnumValueDefinition {
	fn := func(value *EnumValueDefinition) bool {
		return value.Name == name
	}
	return firstEnumValue(l, fn)
}

func firstEnumValue(list EnumValueList, filter valueEnumFilter) *EnumValueDefinition {
	for _, value := range list {
		if filter(value) {
			return value
		}
	}
	return nil
}

