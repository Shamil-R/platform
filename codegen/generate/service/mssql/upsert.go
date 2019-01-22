package mssql

import (
	"context"
)

func Upsert(ctx context.Context, result interface{}) error {
	// TODO: refactoring upsert
	// err := Item(ctx, result)
	// var arg string
	// var f ArgName = func() string {return arg}
	// if err == nil {
	// 	arg = "update"
	// 	return Update(ctx, result, f)
	// } else {
	// 	arg = "create"
	// 	return Create(ctx, result, f)
	// }
	return nil
}
