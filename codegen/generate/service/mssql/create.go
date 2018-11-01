package mssql

import (
	"context"
	"fmt"
	"gitlab/nefco/platform/codegen/schema"

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

func create(tx *sqlx.Tx, v *schema.Value, result interface{}) error {
	for _, child := range v.Children() {
		input := child.Value().Definition().Directives().Input()
		if input != nil && input.IsCreateOneWithout() {
			if err := createOneWithout(tx, child.Value()); err != nil {
				return err
			}
		}
	}
	return nil
}

func createOneWithout(tx *sqlx.Tx, v *schema.Value) error {
	if connect := v.Children().Connect(); connect != nil {
		if err := connectOne(tx, connect.Value()); err != nil {
			return err
		}
	}
	return nil
}

func connectOne(tx *sqlx.Tx, v *schema.Value) error {
	def := v.Definition()
	table := def.Directives().Table()
	fmt.Println(table.Name, table.ArgName(), v.Definition().Name)
	// t := v.Definition().Fields()
	// fmt.Println(t)
	// for _, f := range t {
	// 	fmt.Println(f)
	// }
	for _, child := range v.Children() {
		field := def.Fields().ByName(child.Name)

		if primary := field.Directives().Primary(); primary != nil {

		}
		fmt.Println("$", field.Directives().Field().ArgName())
	}
	return nil
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
