package server

import (
	"fmt"
	"gitlab/nefco/platform/service"
	"gitlab/nefco/platform/service/auth"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/handler"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	Port int `mapstructure:"port"`
}

var DefaultConfig = Config{
	Port: 8080,
}

func Run(v *viper.Viper, exec graphql.ExecutableSchema) error {
	cfg := DefaultConfig
	if err := v.UnmarshalKey("app", &cfg); err != nil {
		return err
	}

	logger := zap.L().Named("server")
	logger.Info("server config", zap.Any("cfg", cfg))

	middlewares := service.Middlewares()

	options := make([]handler.Option, 0, len(middlewares))

	for _, m := range middlewares {
		o, err := m.Middleware(v)
		if err != nil {
			return err
		}
		options = append(options, o)
	}
	// вывод Playground
	http.Handle("/", handler.Playground("Platform", "/query"))

	// аутентификация+авторизация
	http.Handle("/login", auth.MiddlewareLogin())

	// регистрация
	http.Handle("/registration", auth.MiddlewareRegistration(v))

	// регистрация
	http.Handle("/confirm", auth.MiddlewareConfirm())

	// graphql запросы+авторизация
	http.Handle("/query", auth.MiddlewareAuth(handler.GraphQL(exec, options...)))
	//http.Handle("/query", handler.GraphQL(exec, options...))

	return http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil)
}
