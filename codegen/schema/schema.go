package schema

import (
	"gitlab/nefco/platform/codegen/schema/graphql"
	"gitlab/nefco/platform/codegen/template"
	"gitlab/nefco/platform/codegen/tools"
	"io/ioutil"

	"github.com/gobuffalo/packr"

	"github.com/vektah/gqlparser"
	"github.com/vektah/gqlparser/ast"
)

func Generate(cfg Config) error {
	tmpl, err := template.Read(packr.NewBox("./templates"))
	if err != nil {
		return err
	}

	file, err := ioutil.ReadFile(cfg.Source)
	if err != nil {
		return err
	}

	source := []*ast.Source{
		&ast.Source{Name: "directives", Input: graphql.Directives},
		&ast.Source{Name: "schema", Input: string(file)},
	}

	schema, gqlErr := gqlparser.LoadSchema(source...)
	if gqlErr != nil {
		return gqlErr
	}

	if err := tools.ExecuteTemplate(tmpl, schema, cfg.Generate); err != nil {
		return err
	}

	return nil
}
