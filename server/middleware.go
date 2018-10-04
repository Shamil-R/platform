package server

import (
	"github.com/99designs/gqlgen/handler"
	"github.com/spf13/viper"
)

type Middleware interface {
	Middleware(v *viper.Viper) (handler.Option, error)
}

var middlewares []Middleware

func RegisterMiddleware(m Middleware) {
	middlewares = append(middlewares, m)
}
