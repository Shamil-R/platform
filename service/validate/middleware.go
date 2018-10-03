package validate

import (
	"context"
	"fmt"
	"gitlab/nefco/platform/service"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/handler"
	"github.com/spf13/viper"
)

func init() {
	service.Register(&validateService{})
}

type validateService struct{}

func (s validateService) Name() string {
	return "validate"
}

func (s validateService) Middleware(v *viper.Viper) (handler.Option, error) {
	return handler.ResolverMiddleware(middleware(&validate{})), nil
}

func middleware(validate Validate) graphql.FieldMiddleware {
	return func(ctx context.Context, next graphql.Resolver) (interface{}, error) {
		if err := validate.Validate(transform(ctx)); err != nil {
			return nil, err
		}
		return next(ctx)
	}
}

func transform(ctx context.Context) *Data {
	resCtx := graphql.GetResolverContext(ctx)
	fmt.Println(resCtx)
	return &Data{}
}
