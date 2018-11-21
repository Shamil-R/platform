package build

import (
	"context"
	"github.com/joomcode/errorx"
	"gitlab/nefco/platform/codegen/generate/service/mssql/query"
)

func  Order(ctx context.Context, query query.Sort) error {
	order, err := extractArgument(ctx, "orderBy")
	if err != nil && !errorx.IsOfType(err, ArgumentDoesNotExist) {
		return err
	} else if order != nil {
		query.SetOrder(order.Raw)
	}

	return nil
}
