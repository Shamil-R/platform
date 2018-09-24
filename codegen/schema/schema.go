package schema

import (
	"bytes"
	"gitlab/nefco/platform/codegen/file"
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

	buff, err := read(cfg.Source)
	if err != nil {
		return err
	}

	schema, err := parse(buff.String())
	if err != nil {
		return err
	}

	if err := tmpl.Execute(buff, schema); err != nil {
		return err
	}

	if err := file.Write(cfg.Generate, buff); err != nil {
		return err
	}

	return nil
}

func read(src string) (*bytes.Buffer, error) {
	box := packr.NewBox("./graphql")

	buff := &bytes.Buffer{}

	directivesFile := box.Bytes("directives.graphql")

	if _, err := buff.Write(directivesFile); err != nil {
		return nil, err
	}

	if _, err := buff.WriteRune('\n'); err != nil {
		return nil, err
	}

	schemaFile, err := ioutil.ReadFile(src)
	if err != nil {
		return nil, err
	}

	if _, err := buff.Write(schemaFile); err != nil {
		return nil, err
	}

	return buff, nil
}

func parse(input string) (*Schema, error) {
	source := &ast.Source{
		Name:  "schema",
		Input: input,
	}

	schema, gqlErr := gqlparser.LoadSchema(source)
	if gqlErr != nil {
		return nil, gqlErr
	}

	return NewSchema(schema), nil
}
