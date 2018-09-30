package app

import (
	"gitlab/nefco/platform/app/graph"
	"gitlab/nefco/platform/codegen/service"
	"gitlab/nefco/platform/server"
	"net/http"

	"github.com/99designs/gqlgen/handler"
	"github.com/spf13/viper"
)

var prepare func(v *viper.Viper) (http.HandlerFunc, error)

func Run(v *viper.Viper) error {
	if prepare == nil {
		panic("server not implemented")
	}
	h, err := prepare(v)
	if err != nil {
		return err
	}
	return server.Run(v, h)
}

func init() {
	prepare = prepareFunc
}

func prepareFunc(v *viper.Viper) (http.HandlerFunc, error) {
	services := service.Services()

	options := make([]handler.Option, len(services))

	for i, s := range services {
		o, err := s.Init(v)
		if err != nil {
			return nil, err
		}
		options[i] = o
	}

	cfg := graph.Config{
		Resolvers: &Resolver{
			// UserService: NewUserService(),
		},
	}

	return handler.GraphQL(graph.NewExecutableSchema(cfg), options...), nil
}
