package service

import (
	"bytes"
	"fmt"
	"gitlab/nefco/platform/codegen/helper"
	"gitlab/nefco/platform/codegen/schema"
	"gitlab/nefco/platform/service"
	"path"
	"strings"

	"github.com/gobuffalo/packr"
)

type Generator interface {
	Generate(a *schema.Action) (string, error)
}

func generator(name string) (Generator, error) {
	for _, s := range service.Services() {
		if s.Name() == name {
			if g, ok := s.(Generator); ok {
				return g, nil
			}
		}
	}
	return nil, fmt.Errorf("'%s' service not implemented", name)
}

var defaultGenerator = "mssql"

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

	gen, err := generator(defaultGenerator)
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
			content, err := gen.Generate(act)
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

		serviceName := strings.ToLower(def.Name) + "_" + cfg.Service.Filename()

		filename := path.Join(cfg.Service.Dir(), serviceName)

		if err := helper.WriteFile(filename, buf); err != nil {
			return err
		}
	}

	return nil
}
