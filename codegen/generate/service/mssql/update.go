package mssql

import (
	"context"
	"fmt"
	"gitlab/nefco/platform/codegen/schema"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/jmoiron/sqlx"
)

func Update(ctx context.Context, result interface{}) error {
	resCtx := graphql.GetResolverContext(ctx)

	tx, err := Begin(ctx)
	if err != nil {
		return err
	}

	field := &schema.Field{Field: resCtx.Field.Field}

	return update(tx, field, result)
	/* resCtx := graphql.GetResolverContext(ctx)

	tx, err := Begin(ctx)
	if err != nil {
		return err
	}

	queryUpdate := query.Update(resCtx.Field.Field)

	fmt.Println(queryUpdate.Query())

	if _, err := tx.NamedExec(queryUpdate.Query(), queryUpdate.Arg()); err != nil {
		return err
	}

	querySelect := query.Select(resCtx.Field.Field)

	fmt.Println(querySelect.Query())

	stmt, err := tx.PrepareNamed(querySelect.Query())
	if err != nil {
		return err
	}

	if err := stmt.Get(result, querySelect.Arg()); err != nil {
		return err
	}

	return nil */
}

func update(tx *sqlx.Tx, f *schema.Field, result interface{}) error {
	query := new(UpdateQuery)

	data := f.Arguments().Data().Value()

	def := data.Definition()
	table := def.Directives().Table().ArgName()

	query.SetTable(table)

	for _, child := range data.Children() {
		fieldDef := def.Fields().ByName(child.Name)

		col := fieldDef.Directives().Field().ArgName()
		val := child.Value().Conv()

		query.AddValue(col, val)
	}

	where := f.Arguments().Where().Value()

	def = where.Definition()

	for _, child := range where.Children() {
		fieldDef := def.Fields().ByName(child.Name)

		col := fieldDef.Directives().Field().ArgName()
		val := child.Value().Conv()

		query.AddСondition(col, val)
	}

	fmt.Println(query.Query(), query.Arg())

	return nil
}

type UpdateQuery struct {
	Value
}

func (q *UpdateQuery) query() string {
	values := make([]string, 0, len(q.values))
	for _, col := range q.values {
		val := fmt.Sprintf("[%s] = :%s", col, col)
		values = append(values, val)
	}
	return strings.Join(values, ", ")
}

func (q *UpdateQuery) Query() string {
	query := fmt.Sprintf(
		"UPDATE %s SET %s %s",
		q.Table.Query(),
		q.query(),
		q.Condition.Query(),
	)
	return query
}
