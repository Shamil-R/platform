package build

import (
	"context"
	"github.com/joomcode/errorx"
	"gitlab/nefco/platform/codegen/generate/service/mssql/query"
	_schema "gitlab/nefco/platform/service/schema"
	"strconv"
)

func Trasher(ctx context.Context, query query.Trasher) error {
	bwithTrashed := false
	bonlyTrashed := false
	//определение необходимости выводить список с sotdDeleted записями
	withTrashed, err := extractArgument(ctx, "withTrashed")
	if err != nil && !errorx.IsOfType(err, ArgumentDoesNotExist) {
		return err
	} else if withTrashed != nil {
		bwithTrashed, err = strconv.ParseBool(withTrashed.Raw)
		if err != nil {
			return err
		}
	}
	//определение необходимости выводить список только из sotdDeleted записей
	onlyTrashed, err := extractArgument(ctx, "onlyTrashed")
	if err != nil && !errorx.IsOfType(err, ArgumentDoesNotExist) {
		return err
	} else if onlyTrashed != nil {
		bonlyTrashed, err = strconv.ParseBool(onlyTrashed.Raw)
		if err != nil {
			return err
		}
	}

	query.SetTrashed(bwithTrashed, bonlyTrashed)

	// get column name of softDelete
	field, err := ExtractField(ctx)
	if err != nil {
		return err
	}

	objectName := field.Definition().Directives().Object().ArgName()

	schemaCtx := _schema.GetContext(ctx)

	if !schemaCtx.Types().ByName(objectName).Directives().SoftDelete().IsDisable() {
		query.SetTrashedFieldName(schemaCtx.Types().ByName(objectName).Directives().SoftDelete().ArgDeleteField())
	}

	return nil
}

