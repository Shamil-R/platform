package codegen

import (
	"gitlab/nefco/platform/codegen/generate/gqlgen"
	"gitlab/nefco/platform/codegen/generate/schema"
	"gitlab/nefco/platform/codegen/generate/server"
	"gitlab/nefco/platform/codegen/generate/service"

	"github.com/spf13/viper"
)

func Generate(v *viper.Viper) error {
	cfg, err := readConfig(v)
	if err != nil {
		return err
	}
	if err := gqlgen.Generate(cfg.GqlgenConfig()); err != nil {
		return err
	}
	if err := schema.Generate(cfg.SchemaConfig()); err != nil {
		return err
	}
	if err := service.Generate(cfg.ServiceConfig()); err != nil {
		return err
	}
	if err := server.Generate(cfg.ServerConfig()); err != nil {
		return err
	}
	return nil
}

func GenerateGqlgen(v *viper.Viper) error {
	cfg, err := readConfig(v)
	if err != nil {
		return err
	}
	return gqlgen.Generate(cfg.GqlgenConfig())
}

func GenerateSchema(v *viper.Viper) error {
	cfg, err := readConfig(v)
	if err != nil {
		return err
	}
	return schema.Generate(cfg.SchemaConfig())
}

func GenerateService(v *viper.Viper) error {
	cfg, err := readConfig(v)
	if err != nil {
		return err
	}
	return service.Generate(cfg.ServiceConfig())
}

func GenerateServer(v *viper.Viper) error {
	cfg, err := readConfig(v)
	if err != nil {
		return err
	}
	return server.Generate(cfg.ServerConfig())
}
