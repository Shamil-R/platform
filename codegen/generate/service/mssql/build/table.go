package build

import (
	"context"
	"gitlab/nefco/platform/codegen/generate/service/mssql/query"
	"gitlab/nefco/platform/codegen/schema"
)

func TableFromContext(ctx context.Context, query query.Table) error {
	field, err := extractField(ctx)
	if err != nil {
		return err
	}
	query.SetTable(field.Definition().Directives().Table().ArgName())
	return nil
}

func TableFromValue(value *schema.Value, query query.Table) error {
	def := value.Definition()
	table := def.Directives().Table()
	if table == nil {
		return DirectiveDoesNotExist.New("table")
	}
	query.SetTable(table.ArgName())
	return nil
}
