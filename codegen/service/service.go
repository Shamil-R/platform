package service

import (
	"bytes"
	"fmt"
	"gitlab/nefco/platform/codegen/file"
	"gitlab/nefco/platform/codegen/helper"
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

var services []Service

func init() {
	services = []Service{
		mssql.New(),
	}
}

func Services() []Service {
	return services
}

func serviceByName(name string) (Service, error) {
	for _, s := range services {
		if s.Name() == name {
			return s, nil
		}
	}
	return nil, fmt.Errorf("'%s' service not implemented", name)
}

var defaultService = "mssql"

type Config struct {
	Schema  string
	Service helper.File
	Model   helper.File
}

func Generate(cfg Config) error {
	schema, err := schema.Load(cfg.Schema)
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
		Code:   code.New(cfg.Service.Package()),
		Schema: sch,
	}
	data.AddImport("context", "context")
	data.AddImport(cfg.Model.Import(), "model")

	buff := &bytes.Buffer{}

	if err := tmpl.Execute(buff, data); err != nil {
		return err
	}

	if err := file.Write(cfg.Service.Path, buff); err != nil {
		return err
	}

	return nil
}

func generateStruct(cfg Config, sch *schema.Schema) error {
	box := packr.NewBox("./templates")

	s, err := serviceByName(defaultService)
	if err != nil {
		return err
	}

	tmplStruct, err := template.Read("service_struct", box)
	if err != nil {
		return err
	}

	tmplStructFunc, err := template.Read("service_struct_func", box)
	if err != nil {
		return err
	}

	for _, def := range sch.Types().ForAction() {
		buff := &bytes.Buffer{}

		data := &struct {
			*code.Code
			*schema.Definition
		}{
			Code:       code.New(cfg.Service.Package()),
			Definition: def,
		}
		data.AddImport("context", "context")
		data.AddImport(cfg.Model.Import(), "model")
		data.AddImport("github.com/jmoiron/sqlx", "sqlx")

		if err := tmplStruct.Execute(buff, data); err != nil {
			return err
		}

		for _, act := range def.Actions() {
			content, err := s.Generate(act)
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

		serviceName := strings.ToLower(def.Name) + "_" + cfg.Service.Filename()

		filename := path.Join(cfg.Service.Dir(), serviceName)

		if err := file.Write(filename, buff); err != nil {
			return err
		}
	}

	return nil
}
