package templates

import (
	"gitlab/nefco/platform/codegen/tools"
	"text/template"

	"github.com/gobuffalo/packr"
)

func Template() (*template.Template, error) {
	box := packr.NewBox("./")

	service, err := tools.ParseTemplate("service", box.String("service.tpl"))
	if err != nil {
		return nil, err
	}

	return service, nil
}
