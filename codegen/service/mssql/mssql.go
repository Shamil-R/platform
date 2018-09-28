package mssql

import (
	"bytes"
	"gitlab/nefco/platform/codegen/schema"
	"gitlab/nefco/platform/codegen/template"

	"github.com/gobuffalo/packr"
)

func Generate(action string, def *schema.Definition) (string, error) {
	box := packr.NewBox("./templates")

	tmpl, err := template.Read(action, box)
	if err != nil {
		return "", err
	}

	buff := &bytes.Buffer{}

	if err := tmpl.Execute(buff, def); err != nil {
		return "", err
	}

	return buff.String(), nil
}
