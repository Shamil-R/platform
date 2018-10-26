package mssql

import (
	"context"
	"fmt"
	"gitlab/nefco/platform/codegen/generate/service/mssql/query"
	"gitlab/nefco/platform/codegen/schema"

	"github.com/99designs/gqlgen/graphql"
	"github.com/jmoiron/sqlx"
	"github.com/vektah/gqlparser/ast"
)

func Create(ctx context.Context, result interface{}) error {
	resCtx := graphql.GetResolverContext(ctx)

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

	return nil
}

func createOne(ctx context.Context, result interface{}) error {
	tx, err := Begin(ctx)
	if err != nil {
		return err
	}

	resCtx := graphql.GetResolverContext(ctx)

	if data := argument(resCtx.Field.Field, "data"); data != nil {
		for _, child := range data.Children {
			if child.Value.Kind == ast.ObjectValue {

			} else {
				createOneWithout(tx, child)
			}
		}
	}

	return nil
}

func createOneWithout(tx *sqlx.Tx, cv *ast.ChildValue) (int, error) {
	if relType := directive(cv, "relation", "type"); relType != nil {
		for _, relChild := range cv.Value.Children {
			switch relChild.Name {
			case "create":
			case "connect":
			}
		}
	}
	return 0, nil
}

func connectOne(tx *sqlx.Tx, cv *schema.ChildValue) (int, error) {
	// tableName := directive(cv, "table", "name")

	// if tableName == nil {
	// 	return 0, errors.New("no directive table")
	// }

	// for _, child := range tableName.Children {
	// 	fieldName := directive(child, "field", "name")
	// }

	return 0, nil
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
}
