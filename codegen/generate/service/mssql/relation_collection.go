package mssql

import (
	"context"
	"gitlab/nefco/platform/codegen/generate/service/mssql/build"
	_query "gitlab/nefco/platform/codegen/generate/service/mssql/query"
)

//todo deleted_at is null
func RelationCollection(ctx context.Context, objID int, result interface{}) error {
	field, err := build.ExtractField(ctx)
	if err != nil {
		return err
	}

	parent := _query.NewSelect()
	parent.SetTable(field.ObjectDefinition().Directives().Table().ArgName())
	parent.AddСondition(field.Definition().Directives().Relation().ArgOwnerKey(), "eq", objID)
	parent.AddColumn(field.Definition().Directives().Relation().ArgForeignKey(), field.Definition().Directives().Relation().ArgForeignKey())

	logQuery(parent)

	tx, err := Begin(ctx)
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareNamed(parent.Query())
	if err != nil {
		return err
	}

	var childID interface{}
	if err := stmt.Get(&childID, parent.Arg()); err != nil {
		return err
	}

	child := _query.NewSelect()

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

	if err := build.RelationCondition(ctx, child, objID); err != nil {
		return err
	}

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

	return nil
}
