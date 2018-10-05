package role

import (
	"gitlab/nefco/auction/errors"
	"fmt"
)

type AvailableFields struct {
	Object string
	Field string
	UserID int
	Action string
}

type Data struct {
	Table  		string
	Field		string
	Action 		string
}

type Role interface {
	CheckAccess (d []Data) error
}

type role struct{}

func (r role) CheckAccess(d []Data) error {
	for _, elem := range d {
		availableFields, err := r.getAvailableFields(1, elem.Table, elem.Field)
		if err != nil {
			return errors.New("Ошибка получения информации о ваших правах")
		}

		available := r.checkAvailability(availableFields, elem.Action)
		if (!available) {
			return errors.New(fmt.Sprintf("У вас не прав на %s %s в %s", elem.Action, elem.Field, elem.Table))
		}
	}
	return nil
}


func (r role) getAvailableFields(userID int, object string, field string) (*AvailableFields, error) {
	action := "update"
	return &AvailableFields{object, field, userID, action}, nil
}

func (r role) checkAvailability(availableFields *AvailableFields, action string) (bool) {
	var res bool
	if (availableFields != nil && (availableFields.Action == action || action == "read")) {
		res = true
	} else {
		res = false
	}
	return res
}


