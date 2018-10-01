package mssql

import (
	"bytes"
	"context"
	"gitlab/nefco/platform/codegen/schema"
	"gitlab/nefco/platform/codegen/template"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/handler"
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
	return handler.RequestMiddleware(middleware()), nil
}

func (s *mssql) Generate(a *schema.Action) (string, error) {
	return generate(a)
}

func middleware() graphql.RequestMiddleware {
	return func(ctx context.Context, next func(ctx context.Context) []byte) []byte {
		return next(ctx)
	}
}

func generate(a *schema.Action) (string, error) {
	box := packr.NewBox("./templates")

	tmpl, err := template.Read(a.Action, box)
	if err != nil {
		return "", err
	}

	buff := &bytes.Buffer{}

	if err := tmpl.Execute(buff, a.Definition); err != nil {
		return "", err
	}

	return buff.String(), nil
}
