package stub

import (
	"context"
	"gitlab/nefco/platform/codegen/helper"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/handler"
	"github.com/spf13/viper"
)

type stub struct{}

func New() *stub {
	return &stub{}
}

func (s *stub) Name() string {
	return "stub"
}

func (s *stub) Init(v *viper.Viper) (handler.Option, error) {
	return handler.RequestMiddleware(middleware()), nil
}

func (s *stub) Generate(a *helper.Action) (string, error) {
	return `panic("not implemented")`, nil
}

func middleware() graphql.RequestMiddleware {
	return func(ctx context.Context, next func(ctx context.Context) []byte) []byte {
		return next(ctx)
	}
}
