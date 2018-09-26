package service

import (
	"bytes"
	"gitlab/nefco/platform/codegen/file"
	"gitlab/nefco/platform/codegen/schema"
	"gitlab/nefco/platform/codegen/service/code"
	"gitlab/nefco/platform/codegen/service/mssql"
	"gitlab/nefco/platform/codegen/template"
	"path"

	"github.com/gobuffalo/packr"
)

type Config struct {
	SchemaPath     string
	ServiceDir     string
	ServicePackage string
	ModelImport    string
}

type Code struct {
	*code.Code
	Schema *schema.Schema
}

func Generate(cfg Config) error {
	box := packr.NewBox("./templates")

	tmpl, err := template.Read("service", box)
	if err != nil {
		return err
	}

	schema, err := schema.Load(cfg.SchemaPath)
	if err != nil {
		return err
	}

	code := &Code{
		Code: &code.Code{
			PackageName: cfg.ServicePackage,
		},
		Schema: schema,
	}
	code.AddImport("context", "context")
	code.AddImport(cfg.ModelImport, "model")

	buff := &bytes.Buffer{}

	if err := tmpl.Execute(buff, code); err != nil {
		return err
	}

	filename := path.Join(cfg.ServiceDir + "service_gen.go")

	if err := file.Write(filename, buff); err != nil {
		return err
	}

	mssqlCfg := mssql.Config{
		ServiceDir:     cfg.ServiceDir,
		ServicePackage: cfg.ServicePackage,
		ModelImport:    cfg.ModelImport,
	}

	if err := mssql.Generate(mssqlCfg, schema); err != nil {
		return err
	}

	return nil
}
