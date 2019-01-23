package build

import (
	"gitlab/nefco/platform/codegen/generate/service/mssql/query"
	"gitlab/nefco/platform/codegen/schema"
	"time"
)


func DefaultValues(value *schema.Value, q query.Values) (error) {
	if err := SoftDelete(value, q); err != nil {
		return err
	}
	if err := Timestamp(value, q); err != nil {
		return err
	}
	return nil
}

func SoftDelete(value *schema.Value, q query.Values) (error) {
	if softDelete := value.Definition().Directives().SoftDelete(); softDelete != nil && !softDelete.IsDisable() {
		q.AddValue(softDelete.ArgDeleteField(), time.Now())
	}

	return nil
}

func Timestamp(value *schema.Value, q query.Values) (error) {
	timestamp := value.Definition().Directives().Timestamp()
	if timestamp != nil && !timestamp.IsDisable() {
		q.AddValue(timestamp.ArgCreateField(), time.Now())
		q.AddValue(timestamp.ArgUpdateField(), time.Now())
	}

	return nil
}
