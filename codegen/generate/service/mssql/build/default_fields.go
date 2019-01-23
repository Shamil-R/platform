package build

import (
	"gitlab/nefco/platform/codegen/generate/service/mssql/query"
	"gitlab/nefco/platform/codegen/schema"
	"time"
)


func DefaulValues(value *schema.Value, q query.Values) (error) {
	t := time.Now()

	if softDelete := value.Definition().Directives().SoftDelete(); softDelete != nil && !softDelete.IsDisable() {
		q.AddValue(softDelete.ArgDeleteField(), t)
	}
	if timestamp := value.Definition().Directives().Timestamp(); timestamp != nil && !timestamp.IsDisable() {
		q.AddValue(timestamp.ArgCreateField(), t)
		q.AddValue(timestamp.ArgUpdateField(), t)
	}

	return nil
}
