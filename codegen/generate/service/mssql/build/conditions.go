package build

import (
	"errors"
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
			for _, or := range child.Value().Children() {
				ConditionsFromValue(or.Value(), query.And())
			}
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

func PrimaryCondition(ctx context.Context, query query.Conditions,id interface{}) error {
	field, err := ExtractField(ctx)
	if err != nil {
		return err
	}
	selection := field.SelectionSet().Fields()
	sel := selection[0]
	primary := sel.ObjectDefinition().Fields().Primary()
	col := primary.Directives().Field().ArgName()

	query.AddСondition(col, "eq", id)

	return nil
}

func RelationCondition(ctx context.Context, query query.Conditions,id interface{}) error {
	field, err := ExtractField(ctx)
	if err != nil {
		return err
	}
	relation := field.Definition().Directives().Relation()
	if relation == nil {
		return errors.New("relation directive in field does not exist")
	}
	var col string
	if relation.ArgType() == "one_to_many" {
		col = relation.ArgForeignKey()
	} else if relation.ArgType() == "many_to_one" {
		col = relation.ArgOwnerKey()
	} else {
		return errors.New("localKey directive in field does not exist")
	}

	query.AddСondition(col, "eq", id)

	return nil
}
