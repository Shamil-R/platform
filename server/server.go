package server

import (
	"fmt"
	"gitlab/nefco/platform/codegen/service"
	"net/http"

	"github.com/99designs/gqlgen/graphql"

	"github.com/99designs/gqlgen/handler"
	"github.com/spf13/viper"
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

	services := service.Services()

	options := make([]handler.Option, len(services))

	for i, s := range services {
		o, err := s.Init(v)
		if err != nil {
			return err
		}
		options[i] = o
	}

	http.Handle("/", handler.Playground("Platform", "/query"))
	http.Handle("/query", handler.GraphQL(exec, options...))

	return http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil)
}
