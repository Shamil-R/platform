package schema

import (
	"gitlab/nefco/platform/codegen/schema/templates"
	"gitlab/nefco/platform/codegen/tools"
)

func Generate(cfg Config) error {
	tmpl, err := templates.Template()
	if err != nil {
		return err
	}

	schema, err := tools.LoadSchema(cfg.Source)
	if err != nil {
		return err
	}

	if err := tools.ExecuteTemplate(tmpl, schema, cfg.Generate); err != nil {
		return err
	}

	return nil
}
