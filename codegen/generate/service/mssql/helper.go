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
	SelectionDoesNotExist = DoesNotExist.New("selection")
	DirectiveDoesNotExist = DoesNotExist.NewSubtype("directive")
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

func fillTable(ctx context.Context, query query.Table) error {
	field, err := extractField(ctx)
	if err != nil {
		return err
	}

	sels := field.SelectionSet().Fields()
	if len(sels) == 0 {
		return SelectionDoesNotExist
	}

	def := sels[0].ObjectDefinition()

	table := def.Directives().Table()
	if table == nil {
		return DirectiveDoesNotExist.New("table")
	}

	name := table.ArgName()
	if name == nil {
		return ArgumentDoesNotExist.New("name")
	}

	query.SetTable(*name)

	return nil
}

func fillColumns(ctx context.Context, query query.Columns) error {
	field, err := extractField(ctx)
	if err != nil {
		return err
	}

	for _, sel := range field.SelectionSet().Fields() {
		directives := sel.Definition().Directives()
		relation := directives.Relation()
		if relation == nil {
			field := directives.Field().ArgName()
			query.AddColumn(*field, sel.Name)
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
		query.AddСondition(*col, val)
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
		query.AddValue(*col, val)
	}
	return nil
}

func useTable(query query.Table, value *schema.Value) error {
	def := value.Definition()

	table := def.Directives().Table()
	if table == nil {
		return DirectiveDoesNotExist.New("table")
	}

	name := table.ArgName()
	if name == nil {
		return ArgumentDoesNotExist.New("name")
	}

	query.SetTable(*name)

	return nil
}

func useColumns(query query.Columns, value *schema.Value) error {
	def := value.Definition()
	for _, child := range value.Children() {
		fieldDef := def.Fields().ByName(child.Name)
		col := fieldDef.Directives().Field().ArgName()
		query.AddColumn(*col, child.Name)
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
		query.AddСondition(*col, val)
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
	return *col, nil
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
	col := relation.ArgForeignKey()
	return *col, nil
}

func logQuery(query query.Query) {
	fmt.Println(query.Query())
	fmt.Println(query.Arg())
}
