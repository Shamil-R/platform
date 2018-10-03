package service

import (
	"github.com/99designs/gqlgen/handler"
	"github.com/spf13/viper"
)

type Service interface {
	Name() string
}

type Middleware interface {
	Service
	Middleware(v *viper.Viper) (handler.Option, error)
}

var services []Service

func Register(s Service) {
	services = append(services, s)
}

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
