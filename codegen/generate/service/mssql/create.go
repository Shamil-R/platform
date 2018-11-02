package mssql

import (
	"context"
	"errors"
	"fmt"
	"gitlab/nefco/platform/codegen/schema"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/jmoiron/sqlx"
)

func Create(ctx context.Context, result interface{}) error {
	resCtx := graphql.GetResolverContext(ctx)

	tx, err := Begin(ctx)
	if err != nil {
		return err
	}

	value := &schema.Value{Value: resCtx.Field.Field.Arguments.ForName("data").Value}

	return create(tx, value, result)
	/* resCtx := graphql.GetResolverContext(ctx)

	queryInsert := query.Insert(resCtx.Field.Field)

	tx, err := Begin(ctx)
	if err != nil {
		return err
	}

	fmt.Println(queryInsert.Query(), queryInsert.Arg())

	res, err := tx.NamedExec(queryInsert.Query(), queryInsert.Arg())
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	where := map[string]interface{}{
		"id": id,
	}

	querySelect := query.SelectWhere(resCtx.Field.Field, where)

	fmt.Println(querySelect.Query())

	stmt, err := tx.PrepareNamed(querySelect.Query())
	if err != nil {
		return err
	}

	if err := stmt.Get(result, querySelect.Arg()); err != nil {
		return err
	}

	return nil */
}

type columns []string

func (c columns) Add(column string) {
	col := fmt.Sprintf("[%s]", column)
	c = append(c, col)
}

func (c columns) String() string {
	return strings.Join(c, ", ")
}

type insertQuery struct {
	table  string
	values map[string]interface{}
}

func (q *insertQuery) SetTable(table string) {
	q.table = table
}

func (q *insertQuery) AddValue(column string, value interface{}) {
	if q.values == nil {
		q.values = map[string]interface{}{}
	}
	q.values[column] = value
}

func (q *insertQuery) Query() string {
	columns := make([]string, 0, len(q.values))
	values := make([]string, 0, len(q.values))
	for col, _ := range q.values {
		column := fmt.Sprintf("[%s]", col)
		columns = append(columns, column)
		value := fmt.Sprintf(":%s", col)
		values = append(values, value)
	}

	query := fmt.Sprintf(
		"INSERT INTO [%s] (%s) VALUES (%s)",
		q.table,
		strings.Join(columns, ", "),
		strings.Join(values, ", "),
	)
	return query
}

func (q *insertQuery) Arg() map[string]interface{} {
	return q.values
}

func create(tx *sqlx.Tx, v *schema.Value, result interface{}) error {
	query := new(insertQuery)

	def := v.Definition()
	table := def.Directives().Table().ArgName()

	query.SetTable(table)

	for _, child := range v.Children() {
		fieldDef := def.Fields().ByName(child.Name)

		field := fieldDef.Directives().Field().ArgName()

		// fmt.Println(child.Name)
		// for _, d := range child.Directives() {
		// 	fmt.Println("d", d.Name)
		// }
		input := child.Value().Definition().Directives().Input()
		if input != nil && input.IsCreateOneWithout() {
			id, err := createOneWithout(tx, child.Value())
			if err != nil {
				return err
			}
			query.AddValue(field, id)
		} else {
			query.AddValue(field, child.Value().Conv())
		}
	}

	fmt.Println(query.Query(), query.Arg())

	res, err := tx.NamedExec(query.Query(), query.Arg())
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	fmt.Println(id)

	return nil
}

func createOneWithout(tx *sqlx.Tx, v *schema.Value) (int, error) {
	if connect := v.Children().Connect(); connect != nil {
		id, err := connectOne(tx, connect.Value())
		if err != nil {
			return 0, err
		}
		return id, nil
	}
	return 0, errors.New("create one failed")
}

type selectQuery struct {
	table   string
	columns []string
	where   map[string]interface{}
}

func (q *selectQuery) SetTable(table string) {
	q.table = table
}

func (q *selectQuery) AddColumn(column string) {
	q.columns = append(q.columns, column)
}

func (q *selectQuery) AddWhere(column string, value interface{}) {
	if q.where == nil {
		q.where = map[string]interface{}{}
	}
	q.where[column] = value
}

func (q *selectQuery) Query() string {
	columns := make([]string, 0, len(q.columns))
	for _, column := range q.columns {
		col := fmt.Sprintf("[%s]", column)
		columns = append(columns, col)
	}

	var where string
	if len(q.where) > 0 {
		var conds []string
		for column, _ := range q.where {
			cond := fmt.Sprintf("[%s] = :%s", column, column)
			conds = append(conds, cond)
		}
		where = fmt.Sprintf(
			"WHERE %s",
			strings.Join(conds, " AND"),
		)
	}
	query := fmt.Sprintf(
		"SELECT %s FROM [%s] %s",
		strings.Join(columns, ", "),
		q.table,
		where,
	)
	return query
}

func (q *selectQuery) Arg() map[string]interface{} {
	return q.where
}

func connectOne(tx *sqlx.Tx, v *schema.Value) (int, error) {
	def := v.Definition()
	table := def.Directives().Table()

	query := new(selectQuery)
	query.SetTable(table.ArgName())

	for _, child := range v.Children() {
		fieldDef := def.Fields().ByName(child.Name)
		if fieldDef.Directives().Primary() != nil {
			field := fieldDef.Directives().Field()
			query.AddColumn(field.ArgName())
			query.AddWhere(field.ArgName(), child.Value().Conv())
		}
	}

	stmt, err := tx.PrepareNamed(query.Query())
	if err != nil {
		return 0, err
	}

	var id int
	if err := stmt.Get(&id, query.Arg()); err != nil {
		return 0, err
	}

	return id, nil
}

/* func createOne(ctx context.Context, result interface{}) error {
	// tx, err := Begin(ctx)
	// if err != nil {
	// 	return err
	// }

	resCtx := graphql.GetResolverContext(ctx)

	if data := argument(resCtx.Field.Field, "data"); data != nil {
		for _, child := range data.Children {
			if child.Value.Kind == ast.ObjectValue {

			} else {
				// createOneWithout(tx, child)
			}
		}
	}

	return nil
}

func createOneWithout(tx *sqlx.Tx, v *schema.ChildValue) (int, error) {
	// create := v.Children().Create()
	connect := v.Children().Connect()

	if connect != nil {
		connectOne(tx, connect)
	}

	// if relType := directive(cv, "relation", "type"); relType != nil {
	// 	for _, relChild := range cv.Value.Children {
	// 		switch relChild.Name {
	// 		case "create":
	// 		case "connect":
	// 		}
	// 	}
	// }
	return 0, nil
}

func connectOne(tx *sqlx.Tx, v *schema.ChildValue) (int, error) {
	table := v.Directives().Table()

	var (
		sel   []string
		where []string
		arg   map[string]interface{} = make(map[string]interface{})
	)

	for _, child := range v.Children() {
		field := child.Directives().Field()

		sel = append(sel, field.ArgName())

		w := fmt.Sprintf("%s = :%s", field.ArgName(), field.ArgName())
		where = append(where, w)

		arg[field.ArgName()] = child.Value().Raw
	}

	query := fmt.Sprintf(
		"SELECT id FROM %s WHERE %s",
		// strings.Join(sel, ", "),
		table.Name,
		strings.Join(where, " AND"),
	)

	stmt, err := tx.PrepareNamed(query)
	if err != nil {
		return 0, err
	}

	var id int
	if err := stmt.Get(&id, arg); err != nil {
		return 0, err
	}

	return id, nil
}

func argument(field *ast.Field, name string) *ast.Value {
	if arg := field.Arguments.ForName(name); arg != nil {
		return arg.Value
	}
	return nil
}

func dir(child *ast.ChildValue, name string) *ast.Directive {
	return child.Value.Definition.Directives.ForName(name)
}

func directive(child *ast.ChildValue, dirName, argName string) *ast.Value {
	if dir := child.Value.Definition.Directives.ForName(dirName); dir != nil {
		if arg := dir.Arguments.ForName(argName); arg != nil {
			return arg.Value
		}
	}
	return nil
} */
