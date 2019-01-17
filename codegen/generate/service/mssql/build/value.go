package build

import (
	"gitlab/nefco/platform/codegen/generate/service/mssql/query"
	"gitlab/nefco/platform/codegen/schema"
)

func Value(value *schema.Value, query query.Values) error {
	for _, child := range value.Children() {
		fieldDef := value.Definition().Fields().ByName(child.Name)
		field := fieldDef.Directives().Field().ArgName()
		query.AddValue(field, child.Value().Conv())
	}
	return nil
}
