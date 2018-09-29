package schema

import (
	"bytes"
	"gitlab/nefco/platform/codegen/file"
	"gitlab/nefco/platform/codegen/template"
	"io"
	"io/ioutil"

	"github.com/gobuffalo/packr"

	"github.com/vektah/gqlparser"
	"github.com/vektah/gqlparser/ast"
)

type Config struct {
	InputSchemaPath  string
	OutputSchemaPath string
}

func Load(filename string) (*Schema, error) {
	buff := &bytes.Buffer{}

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	if _, err := buff.Write(file); err != nil {
		return nil, err
	}

	schema, err := parse(buff.String())
	if err != nil {
		return nil, err
	}

	return schema, nil
}

func Generate(cfg Config) error {
	box := packr.NewBox("./templates")

	tmpl, err := template.Read("schema", box)
	if err != nil {
		return err
	}

	buff := &bytes.Buffer{}

	if err := read(cfg.InputSchemaPath, buff); err != nil {
		return err
	}

	schema, err := parse(buff.String())
	if err != nil {
		return err
	}

	if err := tmpl.Execute(buff, schema); err != nil {
		return err
	}

	if err := file.Write(cfg.OutputSchemaPath, buff); err != nil {
		return err
	}

	return nil
}

func read(src string, wr io.Writer) error {
	box := packr.NewBox("./graphql")

	directivesFile := box.Bytes("directives.graphql")

	if _, err := wr.Write(directivesFile); err != nil {
		return err
	}

	if _, err := wr.Write([]byte("\n")); err != nil {
		return err
	}

	schemaFile, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	if _, err := wr.Write(schemaFile); err != nil {
		return err
	}

	return nil
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
