package build

import (
	"context"
	"gitlab/nefco/platform/codegen/generate/service/mssql/query"
)

func ColumnsFromSelection(ctx context.Context, q query.Columns) error {
	fields, err := extractFieldsFromSelection(ctx)
	if err != nil {
		return err
	}

	for _, field := range fields {
		directives := field.Definition().Directives()
		if !directives.HasRelation() {
			column := directives.Field().ArgName()
			alias := field.Name
			q.AddColumn(column, alias)
		}
	}

	return nil
}
