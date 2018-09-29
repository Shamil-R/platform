package codegen

import (
	"gitlab/nefco/platform/codegen/gqlgen"
	"gitlab/nefco/platform/codegen/schema"
	"gitlab/nefco/platform/codegen/service"

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
