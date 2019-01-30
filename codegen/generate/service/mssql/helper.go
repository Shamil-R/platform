package mssql

import (
	"context"
	"errors"
	"fmt"
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
)

func extractField(ctx context.Context) (*schema.Field, error) {
	resCtx := graphql.GetResolverContext(ctx)
	if resCtx.Field.Field == nil {
		return nil, FieldDoesNotExist
	}

	field := &schema.Field{Field: resCtx.Field.Field}

	return field, nil
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
