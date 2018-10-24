package role

import (
	"context"
	"errors"
	"fmt"
	"github.com/casbin/casbin"
	"gitlab/nefco/platform/codegen/generate/service/mssql"
	"gitlab/nefco/platform/service/auth"

	genserv "gitlab/nefco/platform/app/service"
	model "gitlab/nefco/platform/app/model"
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
	UserID  		string	`json:"user_id"`
	Table  			string	`json:"object"`
	ActionObject 	*string	`json:"action_object"`
	AllowObject 	*bool	`json:"allow_object"`
	Field 			*string	`json:"field"`
	ActionField 	*string	`json:"action_field"`
	AllowField 		*bool	`json:"allow_field"`
}

func GetObject(ctx context.Context, table string, field string, UserID int) ([]*Object, error) {
	tx, err := mssql.Begin(ctx)
	if err != nil {
		return nil, err
	}

	querySelect := `
		SELECT 
			[oa].[user_id], 
			[oa].[object],
			[aoa].[action] as [action_object],
			[aoa].[allow] as [allow_object],
			[fa].[field],
			[afa].[action] as [action_field],
			[afa].[allow] as [allow_field]
		FROM [object_access] [oa] 
		LEFT JOIN [action_object_access] [aoa] ON [aoa].[object_id] = [oa].[id]
		LEFT JOIN [field_access] [fa] ON [fa].[object_id] = [oa].[id]
		LEFT JOIN [action_field_access] [afa] ON [afa].[field_id] = [fa].[id]
		WHERE [oa].[user_id] = :user_id and [oa].[object] = :object and ([fa].[field] = :field or [fa].[field] is null)
	`

	argSelect := map[string]interface{}{
		"user_id": fmt.Sprintf("%d", UserID),
		"object": table,
		"field": field,
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

	tableAllAct := 0
	tableOneAct := 0
	fieldAllAct := 0
	fieldOneAct := 0



	for _, elem := range obj {

		// приоритет 1 - разрешение на таблицу и на все операции
		if elem.ActionObject == nil && elem.AllowObject != nil &&  *elem.AllowObject == true {
			tableAllAct = 1
		}
		if elem.ActionObject == nil && elem.AllowObject != nil &&  *elem.AllowObject != true {
			tableAllAct = -1
		}

		// приоритет 2 - разрешение на таблицу и на определенную операцию
		if elem.ActionObject != nil && *elem.ActionObject == act.Act && elem.AllowObject != nil &&  *elem.AllowObject == true {
			tableOneAct = 1
		}
		if elem.ActionObject != nil && *elem.ActionObject == act.Act && elem.AllowObject != nil &&  *elem.AllowObject != true {
			tableOneAct = -1
		}

		// приоритет 3 - разрешение на поле и на все операции
		if elem.ActionField == nil && elem.AllowField != nil && *elem.AllowField == true {
			fieldAllAct = 1
		}
		if elem.ActionField == nil && elem.AllowField != nil && *elem.AllowField != true {
			fieldAllAct = -1
		}

		// приоритет 4 - разрешение на поле и на определенную операцию
		if elem.ActionField != nil && *elem.ActionField == act.Act && elem.AllowField != nil &&  *elem.AllowField == true {
			fieldOneAct = 1
		}
		if elem.ActionField != nil && *elem.ActionField == act.Act && elem.AllowField != nil &&  *elem.AllowField != true {
			fieldOneAct = -1
		}
	}

	if fieldOneAct == -1 {
		allowed = false
	} else if fieldOneAct == 1 {
		allowed = true
	} else if fieldAllAct == -1 {
		allowed = false
	} else if fieldAllAct == 1 {
		allowed = true
	} else if tableOneAct == -1 {
		allowed = false
	} else if tableOneAct == 1 {
		allowed = true
	} else if tableAllAct == -1 {
		allowed = false
	} else if tableAllAct == 1 {
		allowed = true
	}

	return allowed, nil
}

func ObjectAccess(ctx context.Context) (*model.ObjectAccess, error) {
	// TODO - передавать параметры в функцию
	temp := "create"
	action := []*string{&temp}
	allow := true
	object := "trip"

	userData := auth.GetContext(ctx)
	userID := userData.ID

	objectAccessQuery := genserv.NewObjectAccessQueryService()
	whereObj := model.ObjectAccessWhereInput{1/*, userID, object*/}
	objects, err := objectAccessQuery.ObjectAccesses(ctx, &whereObj)
	if err != nil {
		return nil, err
	}

	var obj model.ObjectAccess
	if objects == nil {
		objectAccessMutation := genserv.NewObjectAccessMutationService()
		createObj := model.ObjectAccessCreateInput{userID, object}
		obj, err = objectAccessMutation.CreateObjectAccess(ctx, createObj)
		if err != nil {
			return nil, err
		}
	}

	/*whereObj := model.ObjectAccessWhereUniqueInput{1, userID, object}
	createObj := model.ObjectAccessCreateInput{userID, object}
	updateObj := model.ObjectAccessUpdateInput{userID, object}
	obj, err := objectAccessMutation.UpsertObjectAccess(ctx, whereObj, createObj, updateObj)*/

	actionObjectAccessMutation := genserv.NewActionObjectAccessMutationService()
	whereAct := model.ActionObjectAccessWhereUniqueInput{1/*, obj.ID, action*/}
	createAct := model.ActionObjectAccessCreateInput{action, &allow}
	updateAct := model.ActionObjectAccessUpdateInput{action, &allow}
	_, err = actionObjectAccessMutation.UpsertActionObjectAccess(ctx, whereAct, createAct, updateAct)
	if err != nil {
		return nil, err
	}

	return &obj, nil
}

func FieldAccess(ctx context.Context) (*model.FieldAccess, error) {
	// TODO - передавать параметры в функцию
	temp := "create"
	action := []*string{&temp}
	allow := true
	//objectID := 1
	field := "adress"

	fieldAccessQuery := genserv.NewFieldAccessQueryService()
	whereFld := model.FieldAccessWhereInput{1/*, objectID, field*/}
	fields, err := fieldAccessQuery.FieldAccesses(ctx, &whereFld)
	if err != nil {
		return nil, err
	}

	var fld model.FieldAccess
	if fields == nil {
		fieldAccessMutation := genserv.NewFieldAccessMutationService()
		createFld := model.FieldAccessCreateInput{/*, objectID, */field}
		fld, err = fieldAccessMutation.CreateFieldAccess(ctx, createFld)
		if err != nil {
			return nil, err
		}
	}

	/*fieldAccessMutation := genserv.NewFieldAccessMutationService()
	whereFld := model.FieldAccessWhereUniqueInput{1}
	createFld := model.FieldAccessCreateInput{field}
	updateFld := model.FieldAccessUpdateInput{field}
	obj, err := objectAccessMutation.UpsertObjectAccess(ctx, whereObj, createObj, updateObj)
	if err != nil {
		return err
	}*/

	actionFieldAccessMutation := genserv.NewActionFieldAccessMutationService()
	whereAct := model.ActionFieldAccessWhereUniqueInput{1/*, fld.ID, action*/}
	createAct := model.ActionFieldAccessCreateInput{action, &allow}
	updateAct := model.ActionFieldAccessUpdateInput{action, &allow}
	_, err = actionFieldAccessMutation.UpsertActionFieldAccess(ctx, whereAct, createAct, updateAct)
	if err != nil {
		return nil, err
	}

	return &fld, nil
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
