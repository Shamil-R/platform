package extension

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/handler"
	"github.com/spf13/viper"
	"github.com/vektah/gqlparser/ast"
	"strings"
)

type key int
const ctxKey key = 3

type service struct{}

type ExtensionContext struct {
	ActionName string
}

func New() *service {
	return &service{}
}

func (s service) Name() string {
	return "extension"
}

func withContext(ctx context.Context, extensionCtx *ExtensionContext) context.Context {
	return context.WithValue(ctx, ctxKey, extensionCtx)
}

func GetContext(ctx context.Context) *ExtensionContext {
	val := ctx.Value(ctxKey)
	if val == nil {
		return nil
	}
	return val.(*ExtensionContext)
}

func (s service) Middleware(v *viper.Viper, schema *ast.Schema) (handler.Option, error) {
	return handler.ResolverMiddleware(middleware()), nil
}

func middleware() graphql.FieldMiddleware {
	return func(ctx context.Context, next graphql.Resolver) (interface{}, error) {
		resCtx := graphql.GetResolverContext(ctx)
		actionName := getActionName(resCtx.Field.Name)
		data := &ExtensionContext{
			actionName,
		}
		ctx = withContext(ctx, data)
		return next(ctx)
	}
}

func getActionName(fieldName string) string {
	actions := [4]string{"create", "update", "delete", "upsert"}

	for _, action := range actions {
		if strings.Contains(fieldName, action) {
			return action
		}
	}
	return "read"
}
