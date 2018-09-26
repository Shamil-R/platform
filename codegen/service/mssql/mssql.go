package mssql

import (
	"bytes"
	"gitlab/nefco/platform/codegen/file"
	"gitlab/nefco/platform/codegen/schema"
	"gitlab/nefco/platform/codegen/service/code"
	"gitlab/nefco/platform/codegen/template"
	"path"
	"strings"

	"github.com/gobuffalo/packr"
)

type Config struct {
	ServiceDir     string
	ServicePackage string
	ModelImport    string
}

type Code struct {
	*code.Code
	Type *schema.Definition
}

func Generate(cfg Config, schema *schema.Schema) error {
	box := packr.NewBox("./templates")

	tmpl, err := template.Read("service", box)
	if err != nil {
		return err
	}

	for _, t := range schema.ObjectTypes() {
		code := &Code{
			Code: &code.Code{
				PackageName: cfg.ServicePackage,
			},
			Type: t,
		}
		code.AddImport("context", "context")
		code.AddImport(cfg.ModelImport, "model")

		buff := &bytes.Buffer{}

		if err := tmpl.Execute(buff, code); err != nil {
			return err
		}

		serviceName := strings.ToLower(t.Name) + "_service_gen.go"

		filename := path.Join(cfg.ServiceDir, serviceName)

		if err := file.Write(filename, buff); err != nil {
			return err
		}
	}

	return nil
}
