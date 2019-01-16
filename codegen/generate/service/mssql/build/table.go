package build

import (
	"context"
	"gitlab/nefco/platform/codegen/generate/service/mssql/query"
)

func Table(ctx context.Context, query query.Table) error {
	field, err := extractField(ctx)
	if err != nil {
		return err
	}
	query.SetTable(field.Definition().Directives().Table().ArgName())
	return nil
}
