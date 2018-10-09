package role

import (
	"context"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/handler"
	"github.com/spf13/viper"
	"github.com/vektah/gqlparser/ast"
)

type service struct{}

func New() *service {
	return &service{}
}

func (s service) Name() string {
	return "role"
}

func (s service) Middleware(v *viper.Viper) (handler.Option, error) {
	return handler.ResolverMiddleware(middleware(&role{})), nil
}

func middleware(role Role) graphql.FieldMiddleware {
	return func(ctx context.Context, next graphql.Resolver) (interface{}, error) {
		data, err := transform(ctx)
		if err != nil {
			return nil, err
		}

		if err := role.CheckAccess(ctx, data); err != nil {
			return nil, err
		}
		return next(ctx)
	}
}

func transform(ctx context.Context) ([]Data, error) {
	resCtx := graphql.GetResolverContext(ctx)
	// фильтр служебных запросов
	name := resCtx.Field.Name
	if strings.HasPrefix(name, "__") {
		return nil, nil
	}
	// фильтр ответов на запрос
	if resCtx.Object != "Query" && resCtx.Object != "Mutation" {
		return nil, nil
	}

	result := []Data{}

	actionName := getActionName(resCtx.Field.Name)

	var err error
	result, err = walkField(resCtx.Field.SelectionSet, actionName, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func walkField(set ast.SelectionSet, actionName string, data []Data) ([]Data, error) {
	if set != nil {
		for _, sel := range set {
			switch sel := sel.(type) {
			case *ast.Field:
				if sel.SelectionSet == nil {
					data = append(data, Data{sel.ObjectDefinition.Name, sel.Name, actionName})
				} else {
					var err error
					data, err = walkField(sel.SelectionSet, actionName, data)
					if err != nil {
						return nil, err
					}
				}
				// {interface,union}
			case *ast.InlineFragment:
				panic("*ast.InlineFragment not implemented")
			case *ast.FragmentSpread:
				panic("*ast.FragmentSpread not implemented")
			}
		}
	}
	return data, nil
}

//todo
func getActionName(fieldName string) string {
	actions := [3]string{"create", "update", "delete"}

	for _, action := range actions {
		if strings.Contains(fieldName, action) {
			return action
		}
	}
	return "read"
}
