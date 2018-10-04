package schema

import (
	"bytes"
	"gitlab/nefco/platform/codegen/helper"
	"gitlab/nefco/platform/codegen/schema"

	"github.com/gobuffalo/packr"
)

type Config struct {
	In  helper.File
	Out helper.File
}

func Generate(cfg Config) error {
	box := packr.NewBox("./templates")

	buf := bytes.NewBuffer(box.Bytes("directives.graphql"))

	if err := schema.ReadSchema(cfg.In.Path, buf); err != nil {
		return err
	}

	tmpl, err := helper.ReadTemplate("schema", box)
	if err != nil {
		return err
	}

	s, err := schema.ParseSchema(buf.String())
	if err != nil {
		return err
	}

	if err := tmpl.Execute(buf, s); err != nil {
		return err
	}

	if err := helper.WriteFile(cfg.Out.Path, buf); err != nil {
		return err
	}

	return nil
}
