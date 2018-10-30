package log

import (
	"context"
	"fmt"
	"gitlab/nefco/platform/service/auth"
	"gitlab/nefco/platform/service/extension"
	"time"
)

type Data struct {
	Object    string
	CreatedAt time.Time
}

type Log interface {
	Save(context.Context, Data) error
}

type log struct{}

func (r log) Save(ctx context.Context, d Data) error {
	userData := auth.GetContext(ctx)
	action := extension.GetContext(ctx)
	fmt.Println(userData.ID, d.Object, d.CreatedAt, action.ActionName)

	return nil
}
