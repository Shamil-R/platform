package codegen

var DefaultConfig = Config{
	Schema: ConfigSchema{
		Source:    "schema.graphql",
		Transform: "schema.gen.graphql",
	},
	Model: ConfigModel{
		Package: "model",
	},
	Service: ConfigService{
		Package:  "service",
		Filename: "service.gen.go",
	},
}

type Config struct {
	Schema  ConfigSchema  `mapstructure:"schema"`
	Model   ConfigModel   `mapstructure:"model"`
	Service ConfigService `mapstructure:"service"`
}

type ConfigSchema struct {
	Source    string `mapstructure:"source"`
	Transform string `mapstructure:"transform"`
}

type ConfigModel struct {
	Package string `mapstructure:"package"`
}

type ConfigService struct {
	Package  string `mapstructure:"package"`
	Filename string `mapstructure:"filename"`
}
