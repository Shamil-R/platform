package validate

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
	return "validate"
}

func (s service) Middleware(v *viper.Viper, schema *ast.Schema) (handler.Option, error) {
	return handler.ResolverMiddleware(middleware(&validate{})), nil
}

func middleware(validate Validate) graphql.FieldMiddleware {
	return func(ctx context.Context, next graphql.Resolver) (interface{}, error) {
		data, err := transform(ctx)
		if err != nil {
			return nil, err
		}

		if err := validate.Validate(data); err != nil {
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
	if resCtx.Object != "Mutation" {
		return nil, nil
	}

	result := []Data{}

	for _, arg := range resCtx.Field.Arguments {

		defFields := arg.Value.Definition.Fields

		var err error
		result, err = walkDefinition(arg.Value.Children, defFields, result)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func walkDefinition(childs ast.ChildValueList, defFields ast.FieldList, data []Data) ([]Data, error) {
	for _, child := range childs {

		if child.Value.Children != nil {
			var err error
			data, err = walkDefinition(child.Value.Children, defFields, data)
			if err != nil {
				return nil, err
			}
		} else {
			//child.Name - имя поля
			//child.Value.Raw - значение поля

			defField := defFields.ForName(child.Name)

			if defField != nil {
				for _, dir := range defField.Directives {
					if dir.Name == "validate" {
						for _, arg := range dir.Arguments {
							data = append(data, Data{arg.Name, arg.Value.Raw, child.Value.Raw, child.Value.Definition.Name})
						}
					}
				}
			}
		}

	}
	return data, nil
}
