package gqlgen

import (
	"gitlab/nefco/platform/codegen/helper"
	"gitlab/nefco/platform/codegen/schema"

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
		return err
	}

	schema, err := schema.ParseSchema(schemaRaw)
	if err != nil {
		return err
	}

	models := codegen.TypeMap{}

	for _, obj := range schema.Types().Objects() {
		fields := map[string]codegen.TypeMapField{}
		for _, field := range obj.Fields().Objects() {
			fields[field.Name] = codegen.TypeMapField{Resolver: true}
		}
		models[obj.Name] = codegen.TypeMapEntry{Fields: fields}
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
		Models: models,
	}

	if err := codegen.Generate(c); err != nil {
		return err
	}

	return nil
}
