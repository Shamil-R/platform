package mssql

import (
	"bytes"
	"gitlab/nefco/platform/codegen/helper"

	"github.com/jmoiron/sqlx"

	"github.com/99designs/gqlgen/handler"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gobuffalo/packr"
	"github.com/spf13/viper"
)

type mssql struct{}

func New() *mssql {
	return &mssql{}
}

func (s *mssql) Name() string {
	return "mssql"
}

func (s *mssql) Init(v *viper.Viper) (handler.Option, error) {
	db, err := sqlx.Connect("", "")
	if err != nil {
		return nil, err
	}
	return handler.RequestMiddleware(middleware(db)), nil
}

func (s *mssql) Generate(a *helper.Action) (string, error) {
	return generate(a)
}

func generate(a *helper.Action) (string, error) {
	box := packr.NewBox("./templates")

	tmpl, err := helper.ReadTemplate(a.Action, box)
	if err != nil {
		return "", err
	}

	buff := &bytes.Buffer{}

	if err := tmpl.Execute(buff, a.Definition); err != nil {
		return "", err
	}

	return buff.String(), nil
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
