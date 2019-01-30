package mssql

import (
	"context"
	"errors"
	"fmt"
	"gitlab/nefco/platform/codegen/generate/service/mssql/build"
	"gitlab/nefco/platform/codegen/generate/service/mssql/query"
	"gitlab/nefco/platform/codegen/schema"

	"github.com/joomcode/errorx"

	"github.com/99designs/gqlgen/graphql"
)

var (
	Erorrs = errorx.NewNamespace("mssql")

	IllegalState = Erorrs.NewType("illegal_state")

	DoesNotExist = IllegalState.NewSubtype("does_not_exist")

	FieldDoesNotExist     = DoesNotExist.New("field")
	ArgumentDoesNotExist  = DoesNotExist.NewSubtype("argument")
)

func extractField(ctx context.Context) (*schema.Field, error) {
	resCtx := graphql.GetResolverContext(ctx)
	if resCtx.Field.Field == nil {
		return nil, FieldDoesNotExist
	}

	field := &schema.Field{Field: resCtx.Field.Field}

	return field, nil
}

func extractArgument(ctx context.Context, name string) (*schema.Value, error) {
	field, err := extractField(ctx)
	if err != nil {
		return nil, err
	}
	argument := field.Arguments().ByName(name)
	if argument == nil {
		return nil, nil
	}
	return argument.Value(), nil
}

func fillTableCondition(ctx context.Context, query query.Table) error {
	data, err := build.ExtractArgument(ctx, "where")
	if err != nil {
		return err
	}

	if err := build.TableFromInput(ctx, data, query); err != nil {
		return err
	}

	return nil
}

func fillValues(ctx context.Context, query query.Values, f ArgName) error {
	data, err := extractArgument(ctx, f())
	if err != nil {
		return err
	}
	def := data.Definition()
	for _, child := range data.Children() {
		fieldDef := def.Fields().ByName(child.Name)
		col := fieldDef.Directives().Field().ArgName()
		val := child.Value().Conv()
		query.AddValue(col, val)
	}
	return nil
}

func getDefaultValues(ctx context.Context, dirName string, argName string) (string, error) {
	argument, err := extractArgument(ctx, "where")
	if err != nil {
		if errorx.IsOfType(err, ArgumentDoesNotExist) {
			return "", nil
		}
		return "", err
	}

	def := argument.Definition()

	dir := def.Definition.Directives.ForName(dirName)
	arg := dir.Arguments.ForName(argName)

	return arg.Value.Raw, nil
}

func getPrimaryColumn(ctx context.Context) (string, error) {
	field, err := extractField(ctx)
	if err != nil {
		return "", err
	}
	selection := field.SelectionSet().Fields()
	sel := selection[0]
	primary := sel.ObjectDefinition().Fields().Primary()
	col := primary.Directives().Field().ArgName()
	return col, nil
}

func getRelationColumn(ctx context.Context) (string, error) {
	field, err := extractField(ctx)
	if err != nil {
		return "", err
	}
	relation := field.Definition().Directives().Relation()
	if relation == nil {
		return "", errors.New("relation directive in field does not exist")
	}
	var col string
	if relation.ArgType() == "one_to_many" {
		col = relation.ArgForeignKey()
	} else if relation.ArgType() == "many_to_one" {
		col = relation.ArgOwnerKey()
	} else {
		return "", errors.New("localKey directive in field does not exist")
	}
	return col, nil
}

func logQuery(query query.Query) {
	fmt.Println(query.Query())
	fmt.Println(query.Arg())
}
