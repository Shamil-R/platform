package codegen

import (
	"github.com/spf13/viper"
	"gitlab/nefco/platform/codegen/generate/gqlgen"
	"gitlab/nefco/platform/codegen/generate/schema"
	"gitlab/nefco/platform/codegen/generate/server"
	"gitlab/nefco/platform/codegen/generate/service"
)

func Generate(v *viper.Viper) error {
	cfg, err := readConfig(v)
	if err != nil {
		return err
	}
	if err := schema.Generate(cfg.SchemaConfig()); err != nil {
		return err
	}
	if err := gqlgen.Generate(cfg.GqlgenConfig()); err != nil {
		return err
	}
	if err := service.Generate(cfg.ServiceConfig()); err != nil {
		return err
	}
	if err := service.Init(cfg.ServiceConfig(), v); err != nil {
		return err
	}
	if err := server.Generate(cfg.ServerConfig()); err != nil {
		return err
	}
	return nil
}

func GenerateSchema(v *viper.Viper) error {
	cfg, err := readConfig(v)
	if err != nil {
		return err
	}
	return schema.Generate(cfg.SchemaConfig())
}

func GenerateGqlgen(v *viper.Viper) error {
	cfg, err := readConfig(v)
	if err != nil {
		return err
	}
	return gqlgen.Generate(cfg.GqlgenConfig())
}

func GenerateService(v *viper.Viper) error {
	cfg, err := readConfig(v)
	if err != nil {
		return err
	}
	return service.Generate(cfg.ServiceConfig())
}

func GenerateInit(v *viper.Viper) error {
	cfg, err := readConfig(v)
	if err != nil {
		return err
	}
	return service.Init(cfg.ServiceConfig(), v)
}

func GenerateServer(v *viper.Viper) error {
	cfg, err := readConfig(v)
	if err != nil {
		return err
	}
	return server.Generate(cfg.ServerConfig())
}
