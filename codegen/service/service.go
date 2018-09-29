package service

import (
	"bytes"
	"gitlab/nefco/platform/codegen/file"
	"gitlab/nefco/platform/codegen/schema"
	"gitlab/nefco/platform/codegen/service/code"
	"gitlab/nefco/platform/codegen/service/mssql"
	"gitlab/nefco/platform/codegen/template"
	"path"
	"strings"

	"github.com/spf13/viper"

	"github.com/99designs/gqlgen/handler"

	"github.com/gobuffalo/packr"
)

type Service interface {
	Name() string
	Init(v *viper.Viper) (handler.Option, error)
	Generate(a *schema.Action) (string, error)
}

type Config struct {
	SchemaPath     string
	OutputDir      string
	ServiceDir     string
	ServicePackage string
	ModelImport    string
}

func Generate(cfg Config) error {
	schema, err := schema.Load(cfg.SchemaPath)
	if err != nil {
		return err
	}

	if err := generateInterface(cfg, schema); err != nil {
		return err
	}

	if err := generateStruct(cfg, schema); err != nil {
		return err
	}

	return nil
}

func generateInterface(cfg Config, sch *schema.Schema) error {
	box := packr.NewBox("./templates")

	tmpl, err := template.Read("service_interface", box)
	if err != nil {
		return err
	}

	data := &struct {
		*code.Code
		Schema *schema.Schema
	}{
		Code:   code.New(cfg.ServicePackage),
		Schema: sch,
	}
	data.AddImport("context", "context")
	data.AddImport(cfg.ModelImport, "model")

	buff := &bytes.Buffer{}

	if err := tmpl.Execute(buff, data); err != nil {
		return err
	}

	filename := path.Join(cfg.ServiceDir + "service_gen.go")

	if err := file.Write(filename, buff); err != nil {
		return err
	}

	return nil
}

func generateStruct(cfg Config, sch *schema.Schema) error {
	box := packr.NewBox("./templates")

	var service Service = mssql.New()

	tmplStruct, err := template.Read("service_struct", box)
	if err != nil {
		return err
	}

	tmplStructFunc, err := template.Read("service_struct_func", box)
	if err != nil {
		return err
	}

	for _, def := range sch.Types().ForMutation() {
		buff := &bytes.Buffer{}

		data := &struct {
			*code.Code
			*schema.Definition
		}{
			Code:       code.New(cfg.ServicePackage),
			Definition: def,
		}
		data.AddImport("context", "context")
		data.AddImport(cfg.ModelImport, "model")
		data.AddImport("github.com/jmoiron/sqlx", "sqlx")

		if err := tmplStruct.Execute(buff, data); err != nil {
			return err
		}

		for _, act := range def.Actions() {
			content, err := service.Generate(act)
			if err != nil {
				return err
			}

			data := &struct {
				*schema.Action
				Content string
			}{
				Action:  act,
				Content: content,
			}

			if err := tmplStructFunc.Execute(buff, data); err != nil {
				return err
			}
		}

		serviceName := strings.ToLower(def.Name) + "_service_gen.go"

		filename := path.Join(cfg.ServiceDir, serviceName)

		if err := file.Write(filename, buff); err != nil {
			return err
		}
	}

	return nil
}
