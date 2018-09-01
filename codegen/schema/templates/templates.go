package templates

import (
	"gitlab/nefco/platform/codegen/tools"
	"text/template"

	"github.com/gobuffalo/packr"
)

func Template() (*template.Template, error) {
	box := packr.NewBox("./")

	schema, err := tools.ParseTemplate("schema", box.String("schema.tpl"))
	if err != nil {
		return nil, err
	}

	return schema, nil
}
