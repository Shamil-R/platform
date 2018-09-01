package service

import (
	"gitlab/nefco/platform/codegen/build"
	"gitlab/nefco/platform/codegen/service/templates"
	"gitlab/nefco/platform/codegen/tools"
)

func Generate(cfg Config) error {
	tmpl, err := templates.Template()
	if err != nil {
		return err
	}

	schema, err := tools.LoadSchema(cfg.Schema)
	if err != nil {
		return err
	}

	code := &build.Code{
		PackageName: cfg.Package(),
		Schema:      schema,
	}

	if err := tools.ExecuteTemplate(tmpl, code, cfg.Filename); err != nil {
		return err
	}

	return nil
}
