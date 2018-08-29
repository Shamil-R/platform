package codegen

import (
	"gitlab/nefco/platform/codegen/schema"
)

var DefaultConfig = Config{
	SchemaConfig: SchemaConfig{
		Source:   "schema.graphql",
		Generate: "schema_gen.graphql",
	},
	ModelConfig: ConfigModel{
		Package: "model",
	},
	ServiceConfig: ConfigService{
		Package:  "service",
		Filename: "service_gen.go",
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

type SchemaConfig struct {
	Source   string `mapstructure:"source"`
	Generate string `mapstructure:"generate"`
}

type ConfigModel struct {
	Package string `mapstructure:"package"`
}

type ConfigService struct {
	Package  string `mapstructure:"package"`
	Filename string `mapstructure:"filename"`
}
