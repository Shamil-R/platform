package codegen

import "gitlab/nefco/platform/codegen/schema"

var DefaultConfig = Config{
	Schema: schema.DefaultConfig,
	Model: ConfigModel{
		Package: "model",
	},
	Service: ConfigService{
		Package:  "service",
		Filename: "service.gen.go",
	},
}

type Config struct {
	Schema  schema.Config `mapstructure:"schema"`
	Model   ConfigModel   `mapstructure:"model"`
	Service ConfigService `mapstructure:"service"`
}

type ConfigModel struct {
	Package string `mapstructure:"package"`
}

type ConfigService struct {
	Package  string `mapstructure:"package"`
	Filename string `mapstructure:"filename"`
}
