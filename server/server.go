package server

import (
	"context"
	"fmt"
	"gitlab/nefco/platform/server/graph"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/handler"
)

func New() error {
	conf := graph.Config{
		Resolvers: &Resolver{},
	}
	http.Handle("/", handler.Playground("Platform", "/query"))
	http.Handle("/query", handler.GraphQL(
		graph.NewExecutableSchema(conf),
		handler.RequestMiddleware(DBMiddleware(NewDB())),
	))

	return http.ListenAndServe(":8080", nil)
}

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
