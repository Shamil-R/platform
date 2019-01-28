package mssql

import (
	"context"
	"github.com/jmoiron/sqlx"

	"gitlab/nefco/platform/codegen/generate/service/mssql/build"
	"gitlab/nefco/platform/codegen/generate/service/mssql/query"
)

func RelationCollection(ctx context.Context, objID int, result interface{}) error {
	//select parent for get foreignKey value
	//var childID interface{}
	child := query.NewSelect()

	if err := build.TableFromSelection(ctx, child); err != nil {
		return err
	}

	if err := build.ColumnsFromSelection(ctx, child); err != nil {
		return err
	}

	if err := build.Conditions(ctx, child); err != nil {
		return err
	}

	if err := build.Pagination(ctx, child); err != nil {
		return err
	}

	col, err := getRelationColumn(ctx)
	if err != nil {
		return err
	}

	/*field, err := build.ExtractField(ctx)
	if err != nil {
		return err
	}*/

	/*if field.Definition().Directives().Relation().ArgType() == "many_to_one" {
		parent := query.NewSelect()





		parent.SetTable(field.Definition().Directives().Relation().ArgTable())
		parent.AddСondition(field.Definition().Directives().Relation().ArgOwnerKey(), "eq", objID)
		parent.AddColumn(field.Definition().Directives().Relation().ArgForeignKey(), field.Definition().Directives().Relation().ArgForeignKey())

		logQuery(parent)
		fmt.Println("aaaaaaaaaa----")
		fmt.Println(objID)

		tx, err := Begin(ctx)

		if err != nil {
			return err
		}

		stmt, err := tx.PrepareNamed(parent.Query())
		if err != nil {
			return err
		}

		if err := stmt.Get(&childID, parent.Arg()); err != nil {
			return err
		}

		child.AddСondition(col, "eq", childID)

		logQuery(child)

		tx, err = Begin(ctx)
		if err != nil {
			return err
		}

		stmt, err = tx.PrepareNamed(child.Query())
		if err != nil {
			return err
		}

		if err := stmt.Get(result, child.Arg()); err != nil {
			return err
		}


	} else {*/
		//query allUser mutation createUser
		child.AddСondition(col, "eq", objID)

		logQuery(child)

		tx, err := Begin(ctx)
		if err != nil {
			return err
		}

		_query, args, err := sqlx.Named(child.Query(), child.Arg())
		if err != nil {
			return err
		}

		_query, args, err = sqlx.In(_query, args...)
		if err != nil {
			return err
		}
		_query = tx.Rebind(_query)
		if err := tx.Select(result, _query, args...); err != nil {
			return err
		}
	//}

	/*
	//query allMaterials
	child.AddСondition(col, "eq", objID)
	tx, err := Begin(ctx)
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareNamed(child.Query())
	if err != nil {
		return err
	}

	if err := stmt.Get(result, child.Arg()); err != nil {
		return err
	}*/

	return nil
}
