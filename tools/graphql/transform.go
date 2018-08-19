package graphql

import (
	"gitlab/nefco/platform/tools/template"
	"os"
)

func Transform(src, dest string) error {
	schema, err := Load(src)
	if err != nil {
		return err
	}

	tmpl, err := template.Load("../tools/graphql/transform.tpl")
	if err != nil {
		return err
	}

	file, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := tmpl.Execute(file, schema); err != nil {
		return err
	}

	return nil
}
