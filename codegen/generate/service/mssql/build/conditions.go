package build

import (
	"context"
	"gitlab/nefco/platform/codegen/generate/service/mssql/query"
	"gitlab/nefco/platform/codegen/schema"

	"github.com/joomcode/errorx"
)

func Conditions(ctx context.Context, query query.Conditions) error {
	where, err := extractArgument(ctx, "where")
	if err != nil {
		if errorx.IsOfType(err, ArgumentDoesNotExist) {
			return nil
		}
		return err
	}
	return ConditionsFromValue(where, query)
}

func ConditionsFromValue(value *schema.Value, query query.Conditions) error {
	def := value.Definition()

	for _, child := range value.Children() {

		switch child.Name {
		case "AND":
		case "OR":
			for _, or := range child.Value().Children() {
				ConditionsFromValue(or.Value(), query.Or())
			}
		default:
			fieldDef := def.Fields().ByName(child.Name)

			directives := fieldDef.Directives()
			col := directives.Field().ArgName()
			cond := directives.Condition().ArgType()
			val := child.Value().Conv()

			query.AddСondition(col, cond, val)
		}
	}

	return nil
}
