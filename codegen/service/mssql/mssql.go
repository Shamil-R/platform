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
	TypeName string
	Schema   *schema.Schema
}

func Generate(cfg Config, s *schema.Schema) error {
	box := packr.NewBox("./templates")

	tmpl, err := template.Read("service", box)
	if err != nil {
		return err
	}

	for _, def := range s.Types().ForMutation() {
		code := &Code{
			Code: &code.Code{
				PackageName: cfg.ServicePackage,
			},
			TypeName: def.Name,
			Schema:   s,
		}
		code.AddImport("context", "context")
		code.AddImport(cfg.ModelImport, "model")

		buff := &bytes.Buffer{}

		if err := tmpl.Execute(buff, code); err != nil {
			return err
		}

		serviceName := strings.ToLower(def.Name) + "_service_gen.go"

		filename := path.Join(cfg.ServiceDir, serviceName)

		if err := file.Write(filename, buff); err != nil {
			return err
		}
	}

	return nil
}
