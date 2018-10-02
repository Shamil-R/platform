package schema

import (
	"bytes"
	"gitlab/nefco/platform/codegen/helper"

	"github.com/gobuffalo/packr"
)

type Config struct {
	In  helper.File
	Out helper.File
}

func Generate(cfg Config) error {
	box := packr.NewBox("./templates")

	buf := bytes.NewBuffer(box.Bytes("directives.graphql"))

	if err := helper.ReadSchema(cfg.In.Path, buf); err != nil {
		return err
	}

	tmpl, err := helper.ReadTemplate("schema", box)
	if err != nil {
		return err
	}

	schema, err := helper.ParseSchema(buf.String())
	if err != nil {
		return err
	}

	if err := tmpl.Execute(buf, schema); err != nil {
		return err
	}

	if err := helper.WriteFile(cfg.Out.Path, buf); err != nil {
		return err
	}

	return nil
}
