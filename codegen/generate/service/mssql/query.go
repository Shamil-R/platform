package mssql

import (
	"github.com/vektah/gqlparser/ast"
)

func collectPagination(field *ast.Field) ([]string, map[string]interface{}) {
	var query []string
	arg := make(map[string]interface{})

	args := paginationArgs(field)
	if len(args) == 0 {
		return query, arg
	}

	skip := args["skip"]

	if skip != nil {
		cond := "[num] > :skip"
		query = append(query, cond)
		arg["skip"] = skip.Value.Raw
	}

	return query, arg
}

func paginationArgs(field *ast.Field) map[string]*ast.Argument {
	args := make(map[string]*ast.Argument)
	names := []string{"skip", "first", "last"}
	for _, name := range names {
		if arg := field.Arguments.ForName(name); arg != nil {
			args[name] = arg
		}
	}
	return args
}
