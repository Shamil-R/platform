package role

import (
	"context"
	"errors"
	"fmt"
	"github.com/casbin/casbin"
	"gitlab/nefco/platform/codegen/generate/service/mssql"
	"gitlab/nefco/platform/service/auth"
)

type Data struct {
	Table  string
	Field  string
	Action string
}

type Role interface {
	CheckAccess(context.Context, []Data) error
}

type role struct{}

type Subject struct {
	UserID int
}

type Action struct {
	Act  string
}

type Object struct {
	UserID  	string	`json:"user_id"`
	Table  		string	`json:"object"`
	ActionTable *string	`json:"action_object"`
	Field 		*string	`json:"field"`
	ActionField *string	`json:"action_field"`
}

func GetObject(ctx context.Context, table string, fied string, UserID int) ([]*Object, error) {
	tx, err := mssql.Begin(ctx)
	if err != nil {
		return nil, err
	}

	querySelect := `
		SELECT 
			[oa].[user_id], 
			[oa].[object],
			[aoa].[action] as [action_object],
			[fa].[field],
			[afa].[action] as [action_field]
		FROM [object_access] [oa] 
		LEFT JOIN [action_object_access] [aoa] ON [aoa].[object_id] = [oa].[id]
		LEFT JOIN [field_access] [fa] ON [fa].[object_id] = [oa].[id]
		LEFT JOIN [action_field_access] [afa] ON [afa].[field_id] = [fa].[id]
		WHERE [oa].[user_id] = :user_id and [oa].[object] = :object and ([fa].[field] = :field or [fa].[field] is null)
	`

	argSelect := map[string]interface{}{
		"user_id": fmt.Sprintf("%d", UserID),
		"object": table,
		"field": fied,
	}

	stmt, err := tx.PrepareNamed(querySelect)
	if err != nil {
		return nil, err
	}

	resSelect := []*Object{}

	if err := stmt.Select(&resSelect, argSelect); err != nil {
		return nil, err
	}

	return resSelect, nil
}

func checkAvailability(args ...interface{}) (interface{}, error) {
	//sub := args[0].(Subject)
	obj := args[1].([]*Object)
	act := args[2].(Action)

	allowed := false
	for _, elem := range obj {
			// есть запись только в таблице Object, значит разрешены все операции со всеми полями
		if (elem.Field == nil && elem.ActionTable == nil) ||
			// разрешенные типы операций записаны только на таблицу
			(elem.Field == nil && elem.ActionTable != nil && (*elem.ActionTable == act.Act || act.Act == "read")) ||
			// дано разрешение на данное поле в таблице, но не указан тип разрешенной операции, а занчит все операции разрешены
			(elem.Field != nil && elem.ActionField == nil) ||
			// дано разрешение на данное поле в таблице и указан тип разрешенной операции
			(elem.Field != nil && elem.ActionField != nil && (*elem.ActionField == act.Act || act.Act == "read")) {

			allowed = true
		}
	}

	return allowed, nil
}

func (r role) CheckAccess(ctx context.Context, d []Data) error {
	e := casbin.NewEnforcer("service/role/abac_model.conf")
	e.AddFunction("my_func", checkAvailability)

	userData := auth.GetContext(ctx)

	for _, elem := range d {
		act := Action{elem.Action}
		sub := Subject{userData.ID}
		obj, err := GetObject(ctx, elem.Table, elem.Field, userData.ID)
		if err != nil {
			return err
		}

		if e.Enforce(sub, obj, act) != true {
			return errors.New("Вам отказано в доступе к запрашиваемым данным")
		}
	}
	return nil
}
