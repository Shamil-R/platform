package schema

import (
	"github.com/vektah/gqlparser/ast"
)

type EnumFieldDefinition struct {
	*ast.EnumValueDefinition
	directives DirectiveList
}

func (f *EnumFieldDefinition) Directives() DirectiveList {
	if f.directives == nil {
		f.directives = make(DirectiveList, 0, len(f.EnumValueDefinition.Directives))
		for _, directive := range f.EnumValueDefinition.Directives {
			dir := &Directive{Directive: directive}
			f.directives = append(f.directives, dir)
		}
	}
	return f.directives
}

type EnumFieldList []*EnumFieldDefinition

type fieldEnumFilter func(field *EnumFieldDefinition) bool

func (l EnumFieldList) ByName(name string) *EnumFieldDefinition {
	fn := func(field *EnumFieldDefinition) bool {
		return field.Name == name
	}
	return firstEnumField(l, fn)
}

func firstEnumField(list EnumFieldList, filter fieldEnumFilter) *EnumFieldDefinition {
	for _, field := range list {
		if filter(field) {
			return field
		}
	}
	return nil
}

