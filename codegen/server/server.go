package server

import (
	"bytes"
	"gitlab/nefco/platform/codegen/file"
	"gitlab/nefco/platform/codegen/helper"
	"gitlab/nefco/platform/codegen/schema"
	"gitlab/nefco/platform/codegen/template"

	"github.com/gobuffalo/packr"
)

type Config struct {
	Schema  string
	Server  helper.Output
	Exec    helper.Output
	Service helper.Output
}

func Generate(cfg Config) error {
	box := packr.NewBox("./templates")

	tmpl, err := template.Read("server", box)
	if err != nil {
		return err
	}

	s, err := schema.Load(cfg.Schema)
	if err != nil {
		return err
	}

	data := &struct {
		*Config
		Types schema.DefinitionList
	}{
		Config: &cfg,
		Types:  s.Types().ForAction(),
	}

	buff := &bytes.Buffer{}

	if err := tmpl.Execute(buff, data); err != nil {
		return err
	}

	if err := file.Write(cfg.Server.Path, buff); err != nil {
		return err
	}

	return nil
}
