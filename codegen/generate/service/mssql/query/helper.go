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

func argumentChildren(input Input, name string) (bool, ast.ChildValueList) {
	if arg := input.Field().Arguments.ForName(name); arg != nil {
		if len(arg.Value.Children) > 0 {
			return true, arg.Value.Children
		}
	}
	return false, nil
}
