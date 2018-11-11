package mssql

import (
	"context"
	"errors"
	"fmt"
	"gitlab/nefco/platform/codegen/generate/service/mssql/query"
	"gitlab/nefco/platform/codegen/schema"

	"github.com/99designs/gqlgen/graphql"
)

func extractField(ctx context.Context) (*schema.Field, error) {
	resCtx := graphql.GetResolverContext(ctx)
	if resCtx.Field.Field == nil {
		return nil, errors.New("field in context does not exist")
	}
	return &schema.Field{Field: resCtx.Field.Field}, nil
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

func fillTable(ctx context.Context, query query.Table) error {
	field, err := extractField(ctx)
	if err != nil {
		return err
	}
	sel := field.SelectionSet().Fields()[0]
	table := sel.ObjectDefinition().Directives().Table().ArgName()
	query.SetTable(table)
	return nil
}

func fillColumns(ctx context.Context, query query.Columns) error {
	field, err := extractField(ctx)
	if err != nil {
		return err
	}
	for _, sel := range field.SelectionSet().Fields() {
		relation := sel.Definition().Directives().Relation()
		if relation == nil {
			field := sel.Definition().Directives().Field().ArgName()
			query.AddColumn(field)
		} else {
		}
	}
	return nil
}

func fillConditions(ctx context.Context, query query.Conditions) error {
	where, err := extractArgument(ctx, "where")
	if err != nil {
		return err
	}
	if where == nil {
		return nil
	}
	def := where.Definition()
	for _, child := range where.Children() {
		fieldDef := def.Fields().ByName(child.Name)
		col := fieldDef.Directives().Field().ArgName()
		val := child.Value().Conv()
		query.AddСondition(col, val)
	}
	return nil
}

func fillValues(ctx context.Context, query query.Values) error {
	data, err := extractArgument(ctx, "data")
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

func useTable(query query.Table, value *schema.Value) error {
	def := value.Definition()
	table := def.Directives().Table().ArgName()
	query.SetTable(table)
	return nil
}

func useColumns(query query.Columns, value *schema.Value) error {
	def := value.Definition()
	for _, child := range value.Children() {
		fieldDef := def.Fields().ByName(child.Name)
		col := fieldDef.Directives().Field().ArgName()
		query.AddColumn(col)
	}
	return nil
}

func useConditions(query query.Conditions, value *schema.Value) error {
	if value == nil {
		return nil
	}
	def := value.Definition()
	for _, child := range value.Children() {
		fieldDef := def.Fields().ByName(child.Name)
		col := fieldDef.Directives().Field().ArgName()
		val := child.Value().Conv()
		query.AddСondition(col, val)
	}
	return nil
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

func logQuery(query query.Query) {
	fmt.Println(query.Query())
	fmt.Println(query.Arg())
}
