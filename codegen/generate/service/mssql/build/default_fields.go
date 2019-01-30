package build

import (
	"context"
	"database/sql"
	"gitlab/nefco/platform/codegen/generate/service/mssql/query"
	_schema "gitlab/nefco/platform/service/schema"
	"time"
)


func DefaultValues(ctx context.Context, q query.Values) (error) {
	if err := SoftDelete(ctx, q); err != nil {
		return err
	}
	if err := Timestamp(ctx, q); err != nil {
		return err
	}
	return nil
}

func SoftDelete(ctx context.Context, q query.Values) (error) {
	return setSoftDelete(ctx, q, time.Now())
}

func Restore(ctx context.Context, q query.Values) (error) {
	return setSoftDelete(ctx, q, sql.NullString{})
}

func setSoftDelete(ctx context.Context, q query.Values, value interface{}) (error) {
	field, err := ExtractField(ctx)
	if err != nil {
		return err
	}
	objectName := field.Definition().Directives().Object().ArgName()

	schemaCtx := _schema.GetContext(ctx)

	if softDelete := schemaCtx.Types().ByName(objectName).Directives().SoftDelete(); softDelete != nil && !softDelete.IsDisable() {
		q.AddValue(softDelete.ArgDeleteField(), value)
	}

	return nil
}

func Timestamp(ctx context.Context, q query.Values) (error) {
	field, err := ExtractField(ctx)
	if err != nil {
		return err
	}
	objectName := field.Definition().Directives().Object().ArgName()

	schemaCtx := _schema.GetContext(ctx)

	timestamp := schemaCtx.Types().ByName(objectName).Directives().Timestamp()
	if timestamp != nil && !timestamp.IsDisable() {
		q.AddValue(timestamp.ArgCreateField(), time.Now())
		q.AddValue(timestamp.ArgUpdateField(), time.Now())
	}

	return nil
}


func SoftDeleteFromDirective(ctx context.Context, q query.Values) error {
	def, err := extractDefinitionFromSelection(ctx)
	if err != nil {
		return err
	}

	if softDelete := def.Directives().SoftDelete(); softDelete != nil && !softDelete.IsDisable() {
		q.AddValue(softDelete.ArgDeleteField(), time.Now())
	}

	return nil
}

func TimestampFromDirective(ctx context.Context, q query.Values) error {
	def, err := extractDefinitionFromSelection(ctx)
	if err != nil {
		return err
	}
	timestamp := def.Directives().Timestamp()
	if timestamp != nil && !timestamp.IsDisable() {
		q.AddValue(timestamp.ArgCreateField(), time.Now())
		q.AddValue(timestamp.ArgUpdateField(), time.Now())
	}

	return nil
}

func UpdatedFromDirective(ctx context.Context, q query.Values) error {
	def, err := extractDefinitionFromSelection(ctx)
	if err != nil {
		return err
	}
	timestamp := def.Directives().Timestamp()
	if timestamp != nil && !timestamp.IsDisable() {
		q.AddValue(timestamp.ArgUpdateField(), time.Now())
	}

	return nil
}
