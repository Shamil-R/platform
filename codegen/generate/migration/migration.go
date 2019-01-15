package migration

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
	In  helper.File
	Out helper.File
}

func Generate(cfg Config) error {
	box := packr.NewBox("./templates")

	gen, err := generator(defaultGenerator)
	if err != nil {
		return err
	}

	tmpl, err := helper.ReadTemplate("migration", box, gen.Box())
	if err != nil {
		return err
	}

	s, err := schema.LoadSchema(cfg.In.Path)
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



/*

func Generate(cfg Config) error {
	box := packr.NewBox("./templates")

	tmpl, err := helper.ReadTemplate("schema_directive", box)
	if err != nil {
		return err
	}

	s, err := schema.LoadSchema(cfg.In.Path)
	if err != nil {
		return err
	}

	buf := bytes.NewBufferString("")

	if err := tmpl.Execute(buf, s); err != nil {
		return err
	}

	tmpl, err = helper.ReadTemplate("schema_relation", box)
	if err != nil {
		return err
	}

	s, err = schema.ParseSchemaWithDirectives(buf.String())
	if err != nil {
		return err
	}

	buf = bytes.NewBufferString("")

	if err := tmpl.Execute(buf, s); err != nil {
		return err
	}

	tmpl, err = helper.ReadTemplate("schema", box)
	if err != nil {
		return err
	}

	s, err = schema.ParseSchemaWithDirectives(buf.String())
	if err != nil {
		return err
	}

	buf = bytes.NewBufferString("")

	if err := tmpl.Execute(buf, s); err != nil {
		return err
	}

	if err := helper.WriteFile(cfg.Out.Path, buf); err != nil {
		return err
	}

	return nil
}
*/
