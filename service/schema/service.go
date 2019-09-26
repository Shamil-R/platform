package schema

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/handler"
	"github.com/spf13/viper"
	"github.com/vektah/gqlparser/ast"
	"gitlab/nefco/platform/codegen/schema"
)


type service struct {}

type key int

const ctxKey key = 1

func New() *service {
	return &service{}
}

func (s service) Name() string {
	return "schema"
}

func (s service) Middleware(v *viper.Viper, schema *ast.Schema) (handler.Option, error) {
	return handler.RequestMiddleware(middleware(schema)), nil
}

func withContext(ctx context.Context, schemaCtx *schema.Schema) context.Context {
	return context.WithValue(ctx, ctxKey, schemaCtx)
}

func GetContext(ctx context.Context) *schema.Schema {
	val := ctx.Value(ctxKey)
	if val == nil {
		return nil
	}
	return val.(*schema.Schema)
}

func middleware(astSchema *ast.Schema) graphql.RequestMiddleware {
	return func(ctx context.Context, next func(ctx context.Context) []byte) []byte {
		return next(withContext(ctx, &schema.Schema{Schema: astSchema}))
	}
}

