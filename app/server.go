package app

import (
	"gitlab/nefco/platform/server"

	"github.com/99designs/gqlgen/graphql"

	"github.com/spf13/viper"
)

var schema graphql.ExecutableSchema

func Run(v *viper.Viper) error {
	if schema == nil {
		panic("server not implemented")
	}
	return server.Run(v, schema)
}
