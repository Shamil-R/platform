package role

import (
	"context"
	"fmt"
	"errors"
	"gitlab/nefco/platform/service/auth"
)

type AvailableFields struct {
	Object string
	Field  string
	UserID int
	Action string
}

type Data struct {
	Table  string
	Field  string
	Action string
}

type Role interface {
	CheckAccess(context.Context, []Data) error
}

type role struct{}

func (r role) CheckAccess(ctx context.Context, d []Data) error {
	for _, elem := range d {
		availableFields, err := r.getAvailableFields(ctx, elem.Table, elem.Field)
		if err != nil {
			return errors.New("Ошибка получения информации о ваших правах")
		}

		available := r.checkAvailability(availableFields, elem.Action)
		if !available {
			return errors.New(fmt.Sprintf("У вас не прав на %s %s в %s", elem.Action, elem.Field, elem.Table))
		}
	}
	return nil
}

func (r role) getAvailableFields(ctx context.Context, object string, field string) (*[]AvailableFields, error) {
	action := "upsert"
	userData := auth.GetContext(ctx)
	return &[]AvailableFields{{object, field, userData.ID, action}}, nil
}

func (r role) checkAvailability(availableFields *[]AvailableFields, action string) bool {
	res := false
	if availableFields != nil && len(*availableFields) > 0 {
		for _, row := range *availableFields {
			if row.Action == action || action == "read" {
				res = true
				break
			}
		}
	}
	return res
}
