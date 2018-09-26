package codegen

import (
	"gitlab/nefco/platform/codegen/schema"
	"gitlab/nefco/platform/codegen/service"
	"path"
)

var DefaultConfig = Config{
	ProjectPath: "gitlab/nefco/platform/",
	SchemaConfig: SchemaConfig{
		Source:   "server/schema/schema.graphql",
		Generate: "server/schema/schema_gen.graphql",
	},
	ModelConfig: ConfigModel{
		Dir: "server/model/",
	},
	ServiceConfig: ConfigService{
		Dir: "server/service/",
	},
}

type Config struct {
	ProjectPath   string
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
		SchemaPath:     c.SchemaConfig.Generate,
		ServiceDir:     c.ServiceConfig.Dir,
		ServicePackage: path.Base(c.ServiceConfig.Dir),
		ModelImport:    path.Join(c.ProjectPath, c.ModelConfig.Dir),
	}
}

type SchemaConfig struct {
	Source   string `mapstructure:"source"`
	Generate string `mapstructure:"generate"`
}

type ConfigModel struct {
	Dir string `mapstructure:"dir"`
}

type ConfigService struct {
	Dir string `mapstructure:"dir"`
}
