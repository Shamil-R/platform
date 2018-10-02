package server

import (
	"bytes"
	"gitlab/nefco/platform/codegen/helper"

	"github.com/gobuffalo/packr"
)

type Config struct {
	Schema  helper.File
	Server  helper.File
	Exec    helper.File
	Service helper.File
}

func Generate(cfg Config) error {
	box := packr.NewBox("./templates")

	tmpl, err := helper.ReadTemplate("server", box)
	if err != nil {
		return err
	}

	s, err := helper.LoadSchema(cfg.Schema.Path)
	if err != nil {
		return err
	}

	data := &struct {
		*Config
		Types helper.DefinitionList
	}{
		Config: &cfg,
		Types:  s.Types().ForAction(),
	}

	buff := &bytes.Buffer{}

	if err := tmpl.Execute(buff, data); err != nil {
		return err
	}

	if err := helper.WriteFile(cfg.Server.Path, buff); err != nil {
		return err
	}

	return nil
}
