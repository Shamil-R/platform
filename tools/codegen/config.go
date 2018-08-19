package codegen

var DefaultConfig = Config{
	Schema: "schema.graphql",
	Model: ConfigModel{
		Package: "model",
	},
	Service: ConfigService{
		Package:  "service",
		Filename: "service.gen.go",
	},
}

type Config struct {
	Schema  string        `mapstructure:"schema"`
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
