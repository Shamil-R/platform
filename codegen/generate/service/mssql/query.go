package mssql

import (
	"fmt"
	"strings"

	"github.com/vektah/gqlparser/ast"
)

func querySelect(field *ast.Field) string {
	cols := make([]string, 0, len(field.SelectionSet))
	for _, sel := range field.SelectionSet {
		switch sel := sel.(type) {
		case *ast.Field:
			col := fmt.Sprintf("[%s]", sel.Name)
			if ok, val := directiveValue(sel.Definition.Directives, "field", "name"); ok {
				col = fmt.Sprintf("%s AS %s", col, val)
			}
			cols = append(cols, col)
		}
	}
	return strings.Join(cols, ",\n\t")
}

func queryFrom(field *ast.Field) string {
	sel := field.SelectionSet[0]
	selField := sel.(*ast.Field)
	objDef := selField.ObjectDefinition
	name := strings.ToLower(objDef.Name)
	if ok, val := directiveValue(objDef.Directives, "table", "name"); ok {
		name = val
	}
	return fmt.Sprintf("[%s]", name)
}

func queryWhere(arg map[string]interface{}) string {
	if len(arg) == 0 {
		return ""
	}
	conditions := make([]string, 0, len(arg))
	for col, _ := range arg {
		condition := fmt.Sprintf("%s = :%s", col, col)
		conditions = append(conditions, condition)
	}
	return fmt.Sprintf("WHERE %s", strings.Join(conditions, "\nAND"))
}

func qSelect(field *ast.Field) (string, map[string]interface{}) {
	w, a := qWhere(field)

	q := fmt.Sprintf(
		"SELECT %s FROM %s %s",
		qColumns(field),
		qTable(field),
		w,
	)

	return q, a
}

func qColumns(field *ast.Field) string {
	cols := make([]string, 0, len(field.SelectionSet))
	for _, sel := range field.SelectionSet {
		switch sel := sel.(type) {
		case *ast.Field:
			col := fmt.Sprintf("[%s]", sel.Name)
			ok, val := directiveValue(sel.Definition.Directives, "field", "name")
			if ok {
				col = fmt.Sprintf("[%s] AS %s", col, val)
			}
			cols = append(cols, col)
		}
	}
	return strings.Join(cols, ",")
}

func qTable(field *ast.Field) string {
	def := (field.SelectionSet[0].(*ast.Field)).ObjectDefinition

	query := strings.ToLower(def.Name)

	if ok, val := directiveValue(def.Directives, "table", "name"); ok {
		query = val
	}

	query = fmt.Sprintf("[%s]", query)

	if len(paginationArgs(field)) > 0 {
		query = fmt.Sprintf(
			"(SELECT *, ROW_NUMBER() OVER (ORDER BY [id]) AS num FROM %s) AS pagination",
			query,
		)
	}

	return query
}

func qWhere(field *ast.Field) (string, map[string]interface{}) {
	queryWhere, argWhere := collectWhere(field)
	queryPagination, argPagination := collectPagination(field)

	query := append(queryWhere, queryPagination...)

	if len(query) == 0 {
		return "", map[string]interface{}{}
	}

	for k, v := range argPagination {
		argWhere[k] = v
	}

	return "WHERE " + strings.Join(query, " AND"), argWhere
}

func collectWhere(field *ast.Field) ([]string, map[string]interface{}) {
	where := field.Arguments.ForName("where")
	if where == nil {
		return make([]string, 0), make(map[string]interface{}, 0)
	}

	values := where.Value.Children

	query := make([]string, 0, len(values))
	arg := make(map[string]interface{}, len(values))

	for _, val := range values {
		name := val.Name

		cond := fmt.Sprintf("[%s] = :%s", name, name)
		query = append(query, cond)

		arg[name] = val.Value.Raw
	}

	return query, arg
}

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
