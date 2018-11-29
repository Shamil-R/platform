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
	Box() packr.Box
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
	box := packr.NewBox("./templates")

	gen, err := generator(defaultGenerator)
	if err != nil {
		return err
	}

	tmpl, err := helper.ReadTemplate("service", box, gen.Box())
	if err != nil {
		return err
	}

	s, err := schema.LoadSchema(cfg.Schema.Path)
	if err != nil {
		return err
	}

	for _, def := range s.Types().Objects().WholeObjects() {
		buf := bytes.NewBuffer([]byte{})

		data := &struct {
			Config
			*schema.Definition
		}{
			Config:     cfg,
			Definition: def,
		}

		if err := tmpl.Execute(buf, data); err != nil {
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
