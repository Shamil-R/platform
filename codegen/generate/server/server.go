package server

import (
	"bytes"
	"gitlab/nefco/platform/codegen/helper"
	"gitlab/nefco/platform/codegen/schema"

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

	s, err := schema.LoadSchema(cfg.Schema.Path)
	if err != nil {
		return err
	}

	data := &struct {
		*Config
		*schema.Schema
	}{
		Config: &cfg,
		Schema: s,
	}

	buf := bytes.NewBuffer([]byte{})

	if err := tmpl.Execute(buf, data); err != nil {
		return err
	}

	if err := helper.WriteFile(cfg.Server.Path, buf); err != nil {
		return err
	}

	return nil
}
