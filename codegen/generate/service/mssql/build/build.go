package build

import (
	"context"
	"gitlab/nefco/platform/codegen/schema"

	"github.com/99designs/gqlgen/graphql"
	"github.com/joomcode/errorx"
)

var (
	Errors = errorx.NewNamespace("mssql")

	IllegalState = Errors.NewType("illegal_state")

	DoesNotExist = IllegalState.NewSubtype("does_not_exist")

	ArgumentDoesNotExist  = DoesNotExist.NewSubtype("argument")
	DirectiveDoesNotExist = DoesNotExist.NewSubtype("directive")
)

func ExtractField(ctx context.Context) (*schema.Field, error) {
	resCtx := graphql.GetResolverContext(ctx)
	if resCtx.Field.Field == nil {
		return nil, DoesNotExist.New("field does not exist in context")
	}

	field := &schema.Field{Field: resCtx.Field.Field}

	return field, nil
}

func extractFieldsFromSelection(ctx context.Context) (schema.FieldList, error) {
	field, err := ExtractField(ctx)
	if err != nil {
		return nil, err
	}
	return field.SelectionSet().Fields(), nil
}

func extractDefinitionFromSelection(ctx context.Context) (*schema.Definition, error) {
	fields, err := extractFieldsFromSelection(ctx)
	if err != nil {
		return nil, err
	}

	if len(fields) == 0 {
		return nil, DoesNotExist.New("definition does not exist in selection")
	}

	def := fields[0].ObjectDefinition()

	return def, nil
}

func ExtractArgument(ctx context.Context, name string) (*schema.Value, error) {
	return extractArgument(ctx, name)
}

func extractArgument(ctx context.Context, name string) (*schema.Value, error) {
	field, err := ExtractField(ctx)
	if err != nil {
		return nil, err
	}

	arg := field.Arguments().ByName(name)
	if arg == nil {
		err := ArgumentDoesNotExist.New(
			"argument '%s' does not exist in field", name)
		return nil, err
	}

	val := arg.Value()
	if val == nil {
		err := ArgumentDoesNotExist.New(
			"value for argument '%s' does not exist", name)
		return nil, err
	}

	return val, nil
}
