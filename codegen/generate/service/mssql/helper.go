package mssql

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
