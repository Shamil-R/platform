package templates

import (
	"gitlab/nefco/platform/codegen/tools"
	"text/template"

	"github.com/gobuffalo/packr"
)

func Template() (*template.Template, error) {
	box := packr.NewBox("./")

	object, err := tools.ParseTemplate("object", box.String("object.gotpl"))
	if err != nil {
		return nil, err
	}

	enum, err := tools.ParseTemplate("enum", box.String("enum.gotpl"))
	if err != nil {
		return nil, err
	}

	schema, err := tools.ParseTemplate("schema", box.String("schema.gotpl"))
	if err != nil {
		return nil, err
	}

	_, err = schema.AddParseTree(object.Name(), object.Tree)
	if err != nil {
		return nil, err
	}

	_, err = schema.AddParseTree(enum.Name(), enum.Tree)
	if err != nil {
		return nil, err
	}

	return schema, nil
}
