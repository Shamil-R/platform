package service

import (
	"github.com/99designs/gqlgen/handler"
	"github.com/spf13/viper"
	"gitlab/nefco/platform/codegen/generate/service/mssql"
	"gitlab/nefco/platform/codegen/generate/service/stub"
	"gitlab/nefco/platform/service/log"
	"gitlab/nefco/platform/service/role"
	"gitlab/nefco/platform/service/validate"
	"gitlab/nefco/platform/service/extension"
)

func init() {
	services = []Service{
		extension.New(),
		validate.New(),
		role.New(),
		log.New(),
		stub.New(),
		mssql.New(),
	}
}

type Service interface {
	Name() string
}

type Middleware interface {
	Service
	Middleware(v *viper.Viper) (handler.Option, error)
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
