package query

import "github.com/vektah/gqlparser/ast"

func directiveValue(dl ast.DirectiveList, dname, aname string) (bool, string) {
	dir := dl.ForName(dname)
	if dir == nil {
		return false, ""
	}
	arg := dir.Arguments.ForName(aname)
	if arg == nil {
		return false, ""
	}
	val := arg.Value
	if val == nil {
		return false, ""
	}
	return true, val.Raw
}

func argumentValue(input Input, name string) (bool, *ast.Value) {
	arg := input.Field().Arguments.ForName(name)
	if arg != nil && arg.Value != nil {
		return true, arg.Value
	}
	return false, nil
}

func columnName(val *ast.Value, child *ast.ChildValue) string {
	field := val.Definition.Fields.ForName(child.Name)
	if ok, val := directiveValue(field.Directives, "field", "name"); ok {
		return val
	}
	return child.Name
}
