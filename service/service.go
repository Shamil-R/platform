package service

import (
	"github.com/99designs/gqlgen/handler"
	"github.com/spf13/viper"
	"github.com/vektah/gqlparser/ast"
	"gitlab/nefco/platform/codegen/generate/service/mssql"
	"gitlab/nefco/platform/service/schema"
)

func init() {
	services = []Service{
		schema.New(),
		// validate.New(),
		// role.New(),
		// log.New(),
		// stub.New(),
		mssql.New(),
	}
}

type Service interface {
	Name() string
}

type Middleware interface {
	Service
	Middleware(*viper.Viper, *ast.Schema) (handler.Option, error)
}

var services []Service

func Services() []Service {
	return services
}

func Middlewares() []Middleware {
	middlewares := make([]Middleware, 0, len(services))
	for _, s := range services {
		if m, ok := s.(Middleware); ok {
			middlewares = append(middlewares, m)
		}
	}
	return middlewares
}
