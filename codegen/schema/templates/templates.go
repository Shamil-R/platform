//go:generate gorunpkg github.com/gobuffalo/packr/packr

package templates

import (
	"gitlab/nefco/platform/codegen/tools"
	"text/template"

	"github.com/gobuffalo/packr"
)

func Template() (*template.Template, error) {
	box := packr.NewBox("./")

	item, err := tools.ParseTemplate("item", box.String("item.tpl"))
	if err != nil {
		return nil, err
	}

	schema, err := tools.ParseTemplate("schema", box.String("schema.tpl"))
	if err != nil {
		return nil, err
	}

	_, err = schema.AddParseTree(item.Name(), item.Tree)
	if err != nil {
		return nil, err
	}

	return schema, nil
}
