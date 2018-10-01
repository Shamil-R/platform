package server

import (
	"fmt"
	"gitlab/nefco/platform/codegen/service"
	"net/http"

	"github.com/99designs/gqlgen/graphql"

	"github.com/99designs/gqlgen/handler"
	"github.com/spf13/viper"
)

type Config struct {
	Port int `mapstructure:"port"`
}

var DefaultConfig = Config{
	Port: 8080,
}

func Run(v *viper.Viper, exec graphql.ExecutableSchema) error {
	cfg := DefaultConfig
	if err := v.UnmarshalKey("app", &cfg); err != nil {
		return err
	}

	services := service.Services()

	options := make([]handler.Option, len(services))

	for i, s := range services {
		o, err := s.Init(v)
		if err != nil {
			return err
		}
		options[i] = o
	}

	http.Handle("/", handler.Playground("Platform", "/query"))
	http.Handle("/query", handler.GraphQL(exec, options...))

	return http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil)
}

/*
func DBMiddleware(db *DB) graphql.RequestMiddleware {
	return func(ctx context.Context, next func(ctx context.Context) []byte) []byte {
		reqCtx := graphql.GetRequestContext(ctx)

		fmt.Println("RequestMiddleware", reqCtx.Doc.Operations)

		dbCtx := NewDBContext(db)

		ctx = WithDBContext(ctx, dbCtx)

		res := next(ctx)

		fmt.Println("Result", *dbCtx.tx)

		return res
	}
}

type DB struct {
	Connection string
}

func NewDB() *DB {
	return &DB{
		Connection: "connection",
	}
}

type platformContextKey int

const dbContextKey platformContextKey = 1

type DBContext struct {
	DB *DB
	tx *string
}

func NewDBContext(db *DB) *DBContext {
	return &DBContext{
		DB: db,
	}
}

func (c *DBContext) Tx(s string) {
	if c.tx == nil {
		c.tx = &s
	}
}

func WithDBContext(ctx context.Context, db *DBContext) context.Context {
	return context.WithValue(ctx, dbContextKey, db)
}

func GetDBContext(ctx context.Context) *DBContext {
	val := ctx.Value(dbContextKey)
	if val == nil {
		return nil
	}
	return val.(*DBContext)
}
*/
