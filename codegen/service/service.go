package service

import (
	"bytes"
	"gitlab/nefco/platform/codegen/file"
	"gitlab/nefco/platform/codegen/schema"
	"gitlab/nefco/platform/codegen/template"

	"github.com/gobuffalo/packr"
)

func Generate(cfg Config) error {
	box := packr.NewBox("./templates")

	tmpl, err := template.Read("service", box)
	if err != nil {
		return err
	}

	schema, err := schema.Load(cfg.Schema)
	if err != nil {
		return err
	}

	service := NewService(cfg, schema)

	buff := &bytes.Buffer{}

	if err := tmpl.Execute(buff, service); err != nil {
		return err
	}

	if err := file.Write(cfg.Filename, buff); err != nil {
		return err
	}

	return nil
}
