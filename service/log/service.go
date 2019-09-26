package log

import (
	"context"
	"fmt"
	"github.com/vektah/gqlparser/ast"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/handler"
	"github.com/spf13/viper"
)

type service struct{}

func New() *service {
	return &service{}
}

func (s service) Name() string {
	return "log"
}

func (s service) Middleware(v *viper.Viper, schema *ast.Schema) (handler.Option, error) {
	return handler.ResolverMiddleware(middleware(&log{})), nil
}

func middleware(log Log) graphql.FieldMiddleware {
	return func(ctx context.Context, next graphql.Resolver) (interface{}, error) {
		data, err := transform(ctx)
		if err != nil {
			panic(err)
		}

		if err := log.Save(ctx, data); err != nil {
			panic(err)
		}
		return next(ctx)
	}
}

func transform(ctx context.Context) (Data, error) {
	resCtx := graphql.GetResolverContext(ctx)

	result := Data{}

	if resCtx.Parent.Object != "Mutation" {
		return result, nil
	}

	var object string

	for key, arg := range resCtx.Parent.Args {
		object +=  key +  fmt.Sprintf(":%v ", arg)
	}

	return Data{object, time.Now()}, nil
}
