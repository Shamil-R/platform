package codegen

import (
	"gitlab/nefco/platform/codegen/gqlgen"
	"gitlab/nefco/platform/codegen/schema"
	"gitlab/nefco/platform/codegen/server"
	"gitlab/nefco/platform/codegen/service"
	"path"

	"github.com/spf13/viper"
)

const (
	schemaPackage  = "schema"
	execPackage    = "graph"
	modelPackage   = "model"
	servicePackage = "service"
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

func (c Config) withPackage(n, p string) string {
	return path.Join(c.withOutput(n), p)
}

func (c Config) execImport() string {
	return c.withProject(c.withOutput(execPackage))
}

func (c Config) modelImport() string {
	return c.withProject(c.withOutput(modelPackage))
}

func (c Config) serviceImport() string {
	return c.withProject(c.withOutput(servicePackage))
}

func (c Config) schemaGenPath() string {
	return c.withPackage(schemaPackage, "schema_gen.graphql")
}

func (c Config) GqlgenConfig() gqlgen.Config {
	return gqlgen.Config{
		Schema: c.schemaGenPath(),
		Exec: gqlgen.ConfigExec{
			Filename: c.withPackage(execPackage, "graph_gen.go"),
			Package:  execPackage,
		},
		Model: gqlgen.ConfigModel{
			Filename: c.withPackage(modelPackage, "model_gen.go"),
			Package:  modelPackage,
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
		Package:     servicePackage,
		ModelImport: c.modelImport(),
		SchemaPath:  c.schemaGenPath(),
		OutputDir:   c.withOutput(servicePackage),
	}
}

func (c Config) ServerConfig() server.Config {
	return server.Config{
		Package:       "app",
		ExecImport:    c.execImport(),
		ServiceImport: c.serviceImport(),
		SchemaPath:    c.schemaGenPath(),
		OutputPath:    c.withOutput("server_gen.go"),
	}
}

func readConfig(v *viper.Viper) (Config, error) {
	cfg := DefaultConfig
	if err := v.UnmarshalKey("codegen", &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
