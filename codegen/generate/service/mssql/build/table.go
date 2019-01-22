package build

import (
	"context"
	"gitlab/nefco/platform/codegen/generate/service/mssql/query"
	"gitlab/nefco/platform/codegen/schema"
)

func TableFromField(ctx context.Context, q query.Table) error {
	field, err := extractField(ctx)
	if err != nil {
		return err
	}
	q.SetTable(field.Definition().Directives().Table().ArgName())
	return nil
}

func TableFromSelection(ctx context.Context, q query.Table) error {
	def, err := extractDefinitionFromSelection(ctx)
	if err != nil {
		return err
	}
	return tableFromDirective(def, q)
}

func TableFromValue(value *schema.Value, q query.Table) error {
	return tableFromDirective(value.Definition(), q)
}

func tableFromDirective(def *schema.Definition, q query.Table) error {
	table := def.Directives().Table()
	if table == nil {
		return DirectiveDoesNotExist.New("table")
	}
	q.SetTable(table.ArgName())
	return nil
}
