package server

import (
	"fmt"
	"net/http"

	"github.com/99designs/gqlgen/handler"
	"github.com/spf13/viper"
)

type Config struct {
	Port int `mapstructure:"port"`
}

var DefaultConfig = Config{
	Port: 8080,
}

func Run(v *viper.Viper, h http.HandlerFunc) error {
	cfg := DefaultConfig
	if err := v.UnmarshalKey("app", &cfg); err != nil {
		return err
	}
	http.Handle("/", handler.Playground("Platform", "/query"))
	http.Handle("/query", h)
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
