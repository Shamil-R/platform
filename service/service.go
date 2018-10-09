package service

import (
	"github.com/spf13/viper"
	"github.com/99designs/gqlgen/handler"
	"gitlab/nefco/platform/service/validate"
	"gitlab/nefco/platform/service/role"
	"gitlab/nefco/platform/service/log"
)

func init() {
	services = []Service{
		validate.New(),
		role.New(),
		log.New(),
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

func Middlewares()[]Middleware {
	middlewares := make([]Middleware, 0, len(services))
	for _, s := range services {
		if m, ok := s.(Middleware); ok {
			middlewares = append(middlewares, m)
		}
	}
	return middlewares
}
