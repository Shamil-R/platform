package build

import (
	"context"
	"gitlab/nefco/platform/codegen/generate/service/mssql/query"
	"strconv"

	"github.com/joomcode/errorx"
)

func Pagination(ctx context.Context, query query.Pagination) error {
	skip, err := extractArgument(ctx, "skip")
	if err != nil && !errorx.IsOfType(err, ArgumentDoesNotExist) {
		return err
	} else if skip != nil {
		iskip, err := strconv.Atoi(skip.Raw)
		if err != nil {
			return err
		}
		query.SetSkip(iskip)
	}

	first, err := extractArgument(ctx, "first")
	if err != nil && !errorx.IsOfType(err, ArgumentDoesNotExist) {
		return err
	} else if first != nil {
		ifirst, err := strconv.Atoi(first.Raw)
		if err != nil {
			return err
		}
		query.SetFirst(ifirst)
	}

	last, err := extractArgument(ctx, "last")
	if err != nil && !errorx.IsOfType(err, ArgumentDoesNotExist) {
		return err
	} else if last != nil {
		ilast, err := strconv.Atoi(last.Raw)
		if err != nil {
			return err
		}
		query.SetLast(ilast)
	}

	return nil
}
