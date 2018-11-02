package mssql

import (
	"context"
	"fmt"
	"gitlab/nefco/platform/codegen/schema"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/jmoiron/sqlx"
	"github.com/vektah/gqlparser/ast"
)

func dump(sels ast.SelectionSet, s string) {
	for _, sel := range sels {
		switch sel := sel.(type) {
		case *ast.Field:
			fmt.Println(s, sel.Name, sel.ObjectDefinition.Name)
			// for _, dir := range sel.ObjectDefinition.Directives {
			// 	fmt.Println("object-dir", dir.Name)
			// }
			// for _, dir := range sel.Definition.Directives {
			// 	fmt.Println("field-dir", dir.Name)
			// 	for _, arg := range dir.Arguments {
			// 		fmt.Println("field-dir-arg", arg.Name, arg.Value.Raw)
			// 	}
			// }
			dump(sel.SelectionSet, s+"#")
		}
	}
}

func Collection(ctx context.Context, result interface{}) error {
	resCtx := graphql.GetResolverContext(ctx)

	field := &schema.Field{Field: resCtx.Field.Field}

	// fmt.Println(field.Field.ObjectDefinition.Name)

	// dump(field.SelectionSet(), "#")

	// for _, arg := range field.Arguments {
	// 	fmt.Println("arg", arg.Name)
	// }

	tx, err := Begin(ctx)
	if err != nil {
		return err
	}

	// sels := &schema.SelectionSet{SelectionSet: field.Selections}

	if err := selectMany(tx, field); err != nil {
		return err
	}

	/* query := query.Select(resCtx.Field.Field)

	fmt.Println(query.Query(), query.Arg())

	tx, err := Begin(ctx)
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareNamed(query.Query())
	if err != nil {
		return err
	}

	if err := stmt.Select(result, query.Arg()); err != nil {
		return err
	} */

	return nil
}

func selectMany(tx *sqlx.Tx, field *schema.Field) error {
	query := new(Select)

	for i, selField := range field.SelectionSet().Fields() {
		if i == 0 {
			table := selField.ObjectDefinition().Directives().Table().ArgName()
			query.SetTable(table)
		}

		relation := selField.Definition().Directives().Relation()
		if relation != nil {

		} else {
			field := selField.Definition().Directives().Field().ArgName()
			query.AddColumn(field)
		}
		// fmt.Println(selField.Definition().Name)
		// а := field.Definition().Directives().Field().ArgName()
		// query.AddColumn(а)
	}

	fmt.Println(query.Query(), query.Arg())

	return nil
}

type Select struct {
	table      string
	columns    []string
	conditions []string
	arg        map[string]interface{}
}

func (q *Select) SetTable(table string) {
	q.table = table
}

func (q *Select) AddColumn(column string) {
	col := fmt.Sprintf("[%s]", column)
	q.columns = append(q.columns, col)
}

func (q *Select) AddСondition(column string, value interface{}) {
	if q.arg == nil {
		q.arg = map[string]interface{}{}
	}
	q.arg[column] = value
	cond := fmt.Sprintf("[%s] = :%s", column, column)
	q.conditions = append(q.conditions, cond)
}

func (q *Select) Query() string {
	var where string
	if len(q.conditions) > 0 {
		where = fmt.Sprintf(
			"WHERE %s",
			strings.Join(q.conditions, " AND"),
		)
	}
	query := fmt.Sprintf(
		"SELECT %s FROM [%s] %s",
		strings.Join(q.columns, ", "),
		q.table,
		where,
	)
	return query
}

func (q *Select) Arg() map[string]interface{} {
	return q.arg
}
