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

func (c Config) fileSchema() helper.File {
	return helper.File{
		Root: c.Root,
		Path: c.Schema.Path,
	}
}

func (c Config) fileExec() helper.File {
	return helper.File{
		Root: c.Root,
		Path: path.Join(c.Output.Dir, "graph/graph_gen.go"),
	}
}

func (c Config) fileModel() helper.File {
	return helper.File{
		Root: c.Root,
		Path: path.Join(c.Output.Dir, "model/model_gen.go"),
	}
}

func (c Config) fileResolver() helper.File {
	return helper.File{
		Root: c.Root,
		Path: path.Join(c.Output.Dir, "resolver_gen.go"),
		Type: "Resolver",
	}
}

func (c Config) fileExSchema() helper.File {
	return helper.File{
		Root: c.Root,
		Path: path.Join(c.Output.Dir, "schema/schema_gen.graphql"),
	}
}

func (c Config) fileService() helper.File {
	return helper.File{
		Root: c.Root,
		Path: path.Join(c.Output.Dir, "service/service_gen.go"),
	}
}

func (c Config) fileServer() helper.File {
	return helper.File{
		Root: c.Root,
		Path: path.Join(c.Output.Dir, "server_gen.go"),
	}
}

func (c Config) GqlgenConfig() gqlgen.Config {
	return gqlgen.Config{
		Schema:   c.fileExSchema(),
		Exec:     c.fileExec(),
		Model:    c.fileModel(),
		Resolver: c.fileResolver(),
		Dst:      ".gqlgen.yml",
	}
}

func (c Config) SchemaConfig() schema.Config {
	return schema.Config{
		Src: c.Schema.Path,
		Dst: c.fileExSchema().Path,
	}
}

func (c Config) ServiceConfig() service.Config {
	return service.Config{
		Schema:  c.fileExSchema().Path,
		Service: c.fileService(),
		Model:   c.fileModel(),
	}
}

func (c Config) ServerConfig() server.Config {
	return server.Config{
		Schema:  c.fileExSchema().Path,
		Server:  c.fileServer(),
		Exec:    c.fileExec(),
		Service: c.fileService(),
	}
}

func readConfig(v *viper.Viper) (Config, error) {
	cfg := DefaultConfig
	if err := v.UnmarshalKey("codegen", &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
