package schema

import (
	"gitlab/nefco/platform/codegen/template"
	"io/ioutil"

	"github.com/gobuffalo/packr"

	"github.com/vektah/gqlparser"
	"github.com/vektah/gqlparser/ast"
)

func Load(files ...string) (*Schema, error) {
	box := packr.NewBox("./graphql")

	source := make([]*ast.Source, 0, len(files)+1)
	source = append(source,
		&ast.Source{
			Name:  "directives",
			Input: box.String("directives.graphql"),
		},
	)

	for _, filename := range files {
		file, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}

		source = append(source,
			&ast.Source{
				Name:  filename,
				Input: string(file),
			},
		)
	}

	schema, gqlErr := gqlparser.LoadSchema(source...)
	if gqlErr != nil {
		return nil, gqlErr
	}

	return NewSchema(schema), nil
}

func Generate(cfg Config) error {
	box := packr.NewBox("./templates")

	tmpl, err := template.Read("schema", box)
	if err != nil {
		return err
	}

	schema, err := Load(cfg.Source)
	if err != nil {
		return err
	}

	if err := template.Execute(tmpl, schema, cfg.Generate); err != nil {
		return err
	}

	return nil
}
