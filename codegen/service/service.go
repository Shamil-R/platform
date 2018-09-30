package service

import (
	"bytes"
	"fmt"
	"gitlab/nefco/platform/codegen/file"
	"gitlab/nefco/platform/codegen/schema"
	"gitlab/nefco/platform/codegen/service/code"
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
	services = make([]Service, 0)
}

func RegisterService(s Service) {
	services = append(services, s)
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
	Package     string
	ModelImport string
	SchemaPath  string
	OutputDir   string
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
		Code:   code.New(cfg.Package),
		Schema: sch,
	}
	data.AddImport("context", "context")
	data.AddImport(cfg.ModelImport, "model")

	buff := &bytes.Buffer{}

	if err := tmpl.Execute(buff, data); err != nil {
		return err
	}

	filename := path.Join(cfg.OutputDir, "service_gen.go")

	if err := file.Write(filename, buff); err != nil {
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

	for _, def := range sch.Types().ForMutation() {
		buff := &bytes.Buffer{}

		data := &struct {
			*code.Code
			*schema.Definition
		}{
			Code:       code.New(cfg.Package),
			Definition: def,
		}
		data.AddImport("context", "context")
		data.AddImport(cfg.ModelImport, "model")
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

		serviceName := strings.ToLower(def.Name) + "_service_gen.go"

		filename := path.Join(cfg.OutputDir, serviceName)

		if err := file.Write(filename, buff); err != nil {
			return err
		}
	}

	return nil
}
