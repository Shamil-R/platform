package build

import (
	"context"
	"github.com/joomcode/errorx"
	"gitlab/nefco/platform/codegen/generate/service/mssql/query"
	"strconv"
)

func Trasher(ctx context.Context, query query.Trasher) error {
	bwithTrashed := false
	bonlyTrashed := false
	withTrashed, err := extractArgument(ctx, "withTrashed")
	if err != nil && !errorx.IsOfType(err, ArgumentDoesNotExist) {
		return err
	} else if withTrashed != nil {
		bwithTrashed, err = strconv.ParseBool(withTrashed.Raw)
		if err != nil {
			return err
		}
	}

	onlyTrashed, err := extractArgument(ctx, "onlyTrashed")
	if err != nil && !errorx.IsOfType(err, ArgumentDoesNotExist) {
		return err
	} else if onlyTrashed != nil {
		bonlyTrashed, err = strconv.ParseBool(onlyTrashed.Raw)
		if err != nil {
			return err
		}
	}

	query.SetTrashed(bwithTrashed, bonlyTrashed)

	return nil
}

