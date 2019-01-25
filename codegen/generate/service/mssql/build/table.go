package build

import (
	"context"
	"gitlab/nefco/platform/codegen/generate/service/mssql/query"
	"gitlab/nefco/platform/codegen/schema"
	_schema "gitlab/nefco/platform/service/schema"
)

func TableFromSchema(ctx context.Context, q query.Table) error {
	field, err := extractField(ctx)
	if err != nil {
		return err
	}
	objectName := field.Definition().Directives().Object().ArgName()

	schemaCtx := _schema.GetContext(ctx)

	q.SetTable(schemaCtx.Types().ByName(objectName).Directives().Table().ArgName())

	return nil
}

func TableFromSelection(ctx context.Context, q query.Table) error {
	def, err := extractDefinitionFromSelection(ctx)
	if err != nil {
		return err
	}
	return tableFromDirective(def, q)
}

func TableFromInput(ctx context.Context, value *schema.Value, q query.Table) error {
	objectName := value.Definition().Directives().Object().ArgName()
	schemaCtx := _schema.GetContext(ctx)
	return tableFromDirective(schemaCtx.Types().ByName(objectName), q)
}

func tableFromDirective(def *schema.Definition, q query.Table) error {
	table := def.Directives().Table()
	if table == nil {
		return DirectiveDoesNotExist.New("table")
	}
	q.SetTable(table.ArgName())
	return nil
}
