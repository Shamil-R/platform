package log

import (
	"context"
	"github.com/vektah/gqlparser/ast"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/handler"
	"github.com/spf13/viper"
	"strings"
	"time"
)

type service struct{}

func New()*service{
	return &service{}
}

func (s service) Name() string {
	return "log"
}

func (s service) Middleware(v *viper.Viper) (handler.Option, error) {
	return handler.RequestMiddleware(middleware(&log{})), nil
}

func middleware(log Log) graphql.RequestMiddleware {
	return func(ctx context.Context, next func(ctx context.Context) []byte) []byte {
		data, err := transform(ctx)
		if err != nil {panic(err)}

		if err := log.Save(ctx, data); err != nil {panic(err)}
		return next(ctx)
	}
}

func transform(ctx context.Context) ([]Data, error) {
	reqCtx := graphql.GetRequestContext(ctx)

	result := []Data{}

	for _, operation := range reqCtx.Doc.Operations {
		if operation.Operation != "mutation" {continue}

		for _, selection := range operation.SelectionSet {

			name := selection.(*ast.Field).Name
			if strings.HasPrefix(name, "__") {continue}

			actionName := getActionName(name)

			for _, arg := range selection.(*ast.Field).Arguments {
				result = append(result, Data{arg.Value.String(), time.Now(), actionName})
			}
		}
	}
	return result, nil
}

//todo
func getActionName(fieldName string) (string) {
	actions := [3]string{"create", "update", "delete"}

	for _, action := range actions {
		if strings.Contains(fieldName, action) {
			return action
		}
	}
	return "read"
}
