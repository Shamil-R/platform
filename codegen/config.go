package codegen

import (
	"gitlab/nefco/platform/codegen/schema"
	"gitlab/nefco/platform/codegen/service"
)

var DefaultConfig = Config{
	SchemaConfig: SchemaConfig{
		Source:   "server/schema/schema.graphql",
		Generate: "server/schema/schema_gen.graphql",
	},
	ModelConfig: ConfigModel{
		Filename: "server/model/model_gen.go",
	},
	ServiceConfig: ConfigService{
		Filename: "server/service/service_gen.go",
	},
}

type Config struct {
	SchemaConfig  SchemaConfig  `mapstructure:"schema"`
	ModelConfig   ConfigModel   `mapstructure:"model"`
	ServiceConfig ConfigService `mapstructure:"service"`
}

func (c Config) Schema() schema.Config {
	return schema.Config{
		Source:   c.SchemaConfig.Source,
		Generate: c.SchemaConfig.Generate,
	}
}

func (c Config) Service() service.Config {
	return service.Config{
		Filename:      c.ServiceConfig.Filename,
		Schema:        c.SchemaConfig.Generate,
		ModelFilename: c.ModelConfig.Filename,
	}
}

type SchemaConfig struct {
	Source   string `mapstructure:"source"`
	Generate string `mapstructure:"generate"`
}

type ConfigModel struct {
	Filename string `mapstructure:"filename"`
}

type ConfigService struct {
	Filename string `mapstructure:"filename"`
}
