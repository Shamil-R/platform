package codegen

import (
	"gitlab/nefco/platform/codegen/gqlgen"
	"gitlab/nefco/platform/codegen/helper"
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
	Root: "gitlab/nefco/platform/",
	Schema: ConfigSchema{
		Path: "schema.graphql",
	},
	Output: ConfigOutput{
		Dir: "app/",
	},
}

type Config struct {
	Root   string
	Schema ConfigSchema `mapstructure:"schema"`
	Output ConfigOutput `mapstructure:"output"`
}

type ConfigSchema struct {
	Path string `mapstructure:"path"`
}

type ConfigOutput struct {
	Dir string `mapstructure:"dir"`
}

func (c Config) outputSchema() helper.Output {
	return helper.Output{
		Root: c.Root,
		Path: path.Join(c.Output.Dir, "schema/schema_gen.graphql"),
	}
}

func (c Config) outputExec() helper.Output {
	return helper.Output{
		Root: c.Root,
		Path: path.Join(c.Output.Dir, "graph/graph_gen.go"),
	}
}

func (c Config) outputModel() helper.Output {
	return helper.Output{
		Root: c.Root,
		Path: path.Join(c.Output.Dir, "model/model_gen.go"),
	}
}

func (c Config) outputResolver() helper.Output {
	return helper.Output{
		Root: c.Root,
		Path: path.Join(c.Output.Dir, "resolver_gen.go"),
		Type: "Resolver",
	}
}

func (c Config) outputService() helper.Output {
	return helper.Output{
		Root: c.Root,
		Path: path.Join(c.Output.Dir, "service/service_gen.go"),
	}
}

func (c Config) outputServer() helper.Output {
	return helper.Output{
		Root: c.Root,
		Path: path.Join(c.Output.Dir, "server_gen.go"),
	}
}

func (c Config) GqlgenConfig() gqlgen.Config {
	return gqlgen.Config{
		Schema:   c.outputSchema(),
		Exec:     c.outputExec(),
		Model:    c.outputModel(),
		Resolver: c.outputResolver(),
		Dst:      ".gqlgen.yml",
	}
}

func (c Config) SchemaConfig() schema.Config {
	return schema.Config{
		Src: c.Schema.Path,
		Dst: c.outputSchema().Path,
	}
}

func (c Config) ServiceConfig() service.Config {
	return service.Config{
		Schema:  c.outputSchema().Path,
		Service: c.outputService(),
		Model:   c.outputModel(),
	}
}

func (c Config) ServerConfig() server.Config {
	return server.Config{
		Schema:  c.outputSchema().Path,
		Server:  c.outputServer(),
		Exec:    c.outputExec(),
		Service: c.outputService(),
	}
}

func readConfig(v *viper.Viper) (Config, error) {
	cfg := DefaultConfig
	if err := v.UnmarshalKey("codegen", &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
