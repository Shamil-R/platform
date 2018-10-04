package service

import (
	"bytes"
	"fmt"
	"gitlab/nefco/platform/codegen/generate/service/mssql"
	"gitlab/nefco/platform/codegen/helper"
	"gitlab/nefco/platform/codegen/schema"
	"path"
	"strings"

	"github.com/spf13/viper"

	"github.com/99designs/gqlgen/handler"

	"github.com/gobuffalo/packr"
)

type Service interface {
	Name() string
	Init(v *viper.Viper) (handler.Option, error)
	GenerateCommon(d *schema.Definition) (string, error)
	GenerateAction(a *schema.Action) (string, error)
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
	Schema  helper.File
	Service helper.File
	Model   helper.File
}

func Generate(cfg Config) error {
	schema, err := schema.LoadSchema(cfg.Schema.Path)
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

	tmpl, err := helper.ReadTemplate("service_interface", box)
	if err != nil {
		return err
	}

	data := struct {
		Config
		Schema *schema.Schema
	}{
		Config: cfg,
		Schema: sch,
	}

	buff := &bytes.Buffer{}

	if err := tmpl.Execute(buff, &data); err != nil {
		return err
	}

	if err := helper.WriteFile(cfg.Service.Path, buff); err != nil {
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

	tmplStruct, err := helper.ReadTemplate("service_struct", box)
	if err != nil {
		return err
	}

	tmplStructFunc, err := helper.ReadTemplate("service_struct_func", box)
	if err != nil {
		return err
	}

	for _, def := range sch.Types().ForAction() {
		buf := bytes.NewBufferString("")

		data := &struct {
			Config
			*schema.Definition
		}{
			Config:     cfg,
			Definition: def,
		}

		if err := tmplStruct.Execute(buf, data); err != nil {
			return err
		}

		for _, act := range def.Actions() {
			content, err := s.GenerateAction(act)
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

			if err := tmplStructFunc.Execute(buf, data); err != nil {
				return err
			}
		}

		content, err := s.GenerateCommon(def)
		if err != nil {
			return err
		}

		if _, err := buf.WriteString(content); err != nil {
			return err
		}

		serviceName := strings.ToLower(def.Name) + "_" + cfg.Service.Filename()

		filename := path.Join(cfg.Service.Dir(), serviceName)

		if err := helper.WriteFile(filename, buf); err != nil {
			return err
		}
	}

	return nil
}
