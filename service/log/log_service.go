package log

import (
	"context"
	"fmt"
	"gitlab/nefco/platform/service/auth"
	"time"
)

type Data struct {
	Object    string
	CreatedAt time.Time
	Action    string
}

type Log interface {
	Save(context.Context, []Data) error
}

type log struct{}

func (r log) Save(ctx context.Context, d []Data) error {
	userData := auth.GetContext(ctx)
	for _, elem := range d {
		fmt.Println(userData.ID, elem.Object, elem.CreatedAt, elem.Action)
	}
	return nil
}
