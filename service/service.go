package service

import (
	"gitlab/nefco/platform/codegen/generate/service/mssql"
	"gitlab/nefco/platform/service/log"
	"gitlab/nefco/platform/service/role"
	"gitlab/nefco/platform/service/validate"

	"github.com/99designs/gqlgen/handler"
	"github.com/spf13/viper"
)

func init() {
	services = []Service{
		validate.New(),
		role.New(),
		log.New(),
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
