package codegen

import (
	"gitlab/nefco/platform/codegen/gqlgen"
	"gitlab/nefco/platform/codegen/schema"
	"gitlab/nefco/platform/codegen/service"
	"path"

	"github.com/spf13/viper"
)

var DefaultConfig = Config{
	ProjectDir: "gitlab/nefco/platform/",
	Schema: ConfigSchema{
		Path: "schema.graphql",
	},
	Output: ConfigOutput{
		Dir: "app/",
	},
}

type Config struct {
	ProjectDir string
	Schema     ConfigSchema `mapstructure:"schema"`
	Output     ConfigOutput `mapstructure:"output"`
}

type ConfigSchema struct {
	Path string `mapstructure:"path"`
}

type ConfigOutput struct {
	Dir string `mapstructure:"dir"`
}

func (c Config) withProject(p string) string {
	return path.Join(c.ProjectDir, p)
}

func (c Config) withOutput(p string) string {
	return path.Join(c.Output.Dir, p)
}

func (c Config) schemaGenPath() string {
	return c.withOutput("schema/schema_gen.graphql")
}

func (c Config) GqlgenConfig() gqlgen.Config {
	return gqlgen.Config{
		Schema: c.schemaGenPath(),
		Exec: gqlgen.ConfigExec{
			Filename: c.withOutput("graph/graph_gen.go"),
			Package:  "graph",
		},
		Model: gqlgen.ConfigModel{
			Filename: c.withOutput("model/model_gen.go"),
			Package:  "model",
		},
		Resolver: gqlgen.ConfigResolver{
			Filename: c.withOutput("resolver_gen.go"),
			Type:     "Resolver",
		},
		Output: ".gqlgen.yml",
	}
}

func (c Config) SchemaConfig() schema.Config {
	return schema.Config{
		InputSchemaPath:  c.Schema.Path,
		OutputSchemaPath: c.schemaGenPath(),
	}
}

func (c Config) ServiceConfig() service.Config {
	return service.Config{
		Package:     "service",
		ModelImport: c.withProject(c.withOutput("model")),
		SchemaPath:  c.schemaGenPath(),
		OutputDir:   c.withOutput("service/"),
	}
}

func readConfig(v *viper.Viper) (Config, error) {
	cfg := DefaultConfig
	if err := v.UnmarshalKey("codegen", &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
