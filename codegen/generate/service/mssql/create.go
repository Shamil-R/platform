package mssql

import (
	"context"
	"errors"
	"gitlab/nefco/platform/codegen/generate/service/mssql/build"
	"gitlab/nefco/platform/codegen/generate/service/mssql/query"
	"gitlab/nefco/platform/codegen/schema"
)

func Create(ctx context.Context, result interface{}, f ArgName) error {
	if err := create(ctx, result, f); err != nil {
		return err
	}
	return nil
}

func create(ctx context.Context, result interface{}, f ArgName) error {
	query := query.NewInsert()

	if err := fillTable(ctx, query); err != nil {
		return err
	}

	data, err := extractArgument(ctx, f())
	if err != nil {
		return err
	}

	for _, child := range data.Children() {
		fieldDef := data.Definition().Fields().ByName(child.Name)

		field := fieldDef.Directives().Field().ArgName()

		input := child.Value().Definition().Directives().Input()
		if input != nil && input.IsCreateOneWithout() {
			id, err := createOneWithout(ctx, child.Value())
			if err != nil {
				return err
			}
			query.AddValue(field, id)
		} else {
			query.AddValue(field, child.Value().Conv())
		}
	}

	logQuery(query)

	tx, err := Begin(ctx)
	if err != nil {
		return err
	}

	res, err := tx.NamedExec(query.Query(), query.Arg())
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	if err := createResult(ctx, id, result); err != nil {
		return err
	}

	return nil
}

func createOneWithout(ctx context.Context, v *schema.Value) (int64, error) {
	if connect := v.Children().Connect(); connect != nil {
		id, err := connectOne(ctx, connect.Value())
		if err != nil {
			return 0, err
		}
		return id, nil
	}
	return 0, errors.New("create one failed")
}

func connectOne(ctx context.Context, v *schema.Value) (int64, error) {
	query := query.NewSelect()

	if err := useTable(query, v); err != nil {
		return 0, err
	}

	if err := useColumns(query, v); err != nil {
		return 0, err
	}

	if err := build.ConditionsFromValue(v, query); err != nil {
		return 0, err
	}

	logQuery(query)

	tx, err := Begin(ctx)
	if err != nil {
		return 0, err
	}

	stmt, err := tx.PrepareNamed(query.Query())
	if err != nil {
		return 0, err
	}

	var id int64
	if err := stmt.Get(&id, query.Arg()); err != nil {
		return 0, err
	}

	return id, nil
}

func createResult(ctx context.Context, id int64, result interface{}) error {
	query := query.NewSelect()

	if err := fillTable(ctx, query); err != nil {
		return err
	}

	if err := fillColumns(ctx, query); err != nil {
		return err
	}

	col, err := getPrimaryColumn(ctx)
	if err != nil {
		return err
	}

	query.AddСondition(col, "eq", id)

	logQuery(query)

	tx, err := Begin(ctx)
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareNamed(query.Query())
	if err != nil {
		return err
	}

	if err := stmt.Get(result, query.Arg()); err != nil {
		return err
	}

	return nil
}
