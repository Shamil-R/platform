package mssql

import (
	"context"
	"fmt"
	"gitlab/nefco/platform/codegen/generate/service/mssql/query"

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

	if data := argumentValue(resCtx.Field.Field, "data"); data != nil {
		for _, child := range data.Children {
			if child.Value.Kind == ast.ObjectValue {

			} else {
				createOneWithout(tx, child)
			}
		}
	}

	return nil
}

func createOneWithout(tx *sqlx.Tx, child *ast.ChildValue) (int, error) {
	if relType := dirArg(child, "relation", "type"); relType != nil {
		for _, relChild := range child.Value.Children {
			switch relChild.Name {
			case "create":
			case "connect":
			}
		}
	}
	return 0, nil
}

func argumentValue(field *ast.Field, name string) *ast.Value {
	if arg := field.Arguments.ForName(name); arg != nil {
		return arg.Value
	}
	return nil
}

func dirArg(child *ast.ChildValue, dirName, argName string) *ast.Value {
	if dir := child.Value.Definition.Directives.ForName(dirName); dir != nil {
		if arg := dir.Arguments.ForName(argName); arg != nil {
			return arg.Value
		}
	}
	return nil
}
