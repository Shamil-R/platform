package gqlgen

import (
	"fmt"
	"gitlab/nefco/platform/codegen/helper"
	"gitlab/nefco/platform/codegen/schema"
	"os"

	"github.com/99designs/gqlgen/codegen"
)

type Config struct {
	Schema   helper.File
	Exec     helper.File
	Model    helper.File
	Resolver helper.File
}

func Generate(cfg Config) error {
	schemaRaw, err := schema.LoadSchemaRaw(cfg.Schema.Path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "unable to open schema: "+err.Error())
		os.Exit(1)
	}

	c := codegen.Config{
		SchemaStr: schemaRaw,
		Exec: codegen.PackageConfig{
			Filename: cfg.Exec.Path,
			Package:  cfg.Exec.Package(),
		},
		Model: codegen.PackageConfig{
			Filename: cfg.Model.Path,
			Package:  cfg.Model.Package(),
		},
		Resolver: codegen.PackageConfig{
			Filename: cfg.Resolver.Path,
			Type:     cfg.Resolver.Type,
		},
		Models: codegen.TypeMap{
			"User": codegen.TypeMapEntry{
				Fields: map[string]codegen.TypeMapField{
					"materials": codegen.TypeMapField{
						Resolver: true,
					},
				},
			},
			"Material": codegen.TypeMapEntry{
				Fields: map[string]codegen.TypeMapField{
					"author": codegen.TypeMapField{
						Resolver: true,
					},
				},
			},
		},
	}

	if err := codegen.Generate(c); err != nil {
		return err
	}

	return nil
}
