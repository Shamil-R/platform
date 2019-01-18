package build

import (
	"gitlab/nefco/platform/codegen/generate/service/mssql/query"
	"gitlab/nefco/platform/codegen/schema"
)

func Value(value *schema.Value, query query.Values) error {
	for _, child := range value.Children() {
		fieldDef := value.Definition().Fields().ByName(child.Name)
		if !fieldDef.Directives().HasRelation() {
			field := fieldDef.Directives().Field()
			if field != nil {
				query.AddValue(field.ArgName(), child.Value().Conv())
			}
		}
	}
	return nil
}
