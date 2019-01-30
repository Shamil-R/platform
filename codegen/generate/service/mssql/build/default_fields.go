package build

import (
	"context"
	"database/sql"
	"gitlab/nefco/platform/codegen/generate/service/mssql/query"
	"gitlab/nefco/platform/codegen/schema"
	_schema "gitlab/nefco/platform/service/schema"
	"time"
)


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
	timestamp, err := getTimestamp(ctx)
	if err != nil {
		return err
	}
	if timestamp != nil && !timestamp.IsDisable() {
		q.AddValue(timestamp.ArgCreateField(), time.Now())
		q.AddValue(timestamp.ArgUpdateField(), time.Now())
	}

	return nil
}

func Updated(ctx context.Context, q query.Values) (error) {
	timestamp, err := getTimestamp(ctx)
	if err != nil {
		return err
	}
	if timestamp != nil && !timestamp.IsDisable() {
		q.AddValue(timestamp.ArgUpdateField(), time.Now())
	}

	return nil
}

func getTimestamp(ctx context.Context) (*schema.TimestampDirective, error) {
	field, err := ExtractField(ctx)
	if err != nil {
		return nil, err
	}
	objectName := field.Definition().Directives().Object().ArgName()

	schemaCtx := _schema.GetContext(ctx)

	timestamp := schemaCtx.Types().ByName(objectName).Directives().Timestamp()
	return timestamp, nil
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
