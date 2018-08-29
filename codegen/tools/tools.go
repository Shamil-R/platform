package tools

import (
	"gitlab/nefco/platform/codegen/build"
	"io/ioutil"
	"os"
	"path"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/jinzhu/inflection"
	"github.com/vektah/gqlparser"
	"github.com/vektah/gqlparser/ast"
)

func LoadSchema(filename string) (*build.Schema, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	source := &ast.Source{
		Name:  "schema",
		Input: string(file),
	}

	schema, gqlErr := gqlparser.LoadSchema(source)
	if gqlErr != nil {
		return nil, gqlErr
	}

	return build.NewSchema(schema), nil
}

func ParseTemplate(name, text string) (*template.Template, error) {
	funcs := sprig.TxtFuncMap()
	funcs["plural"] = inflection.Plural

	tmpl, err := template.New(name).Funcs(funcs).Parse(text)
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

func ExecuteTemplate(tmpl *template.Template, data interface{}, filename string) error {
	dir := path.Dir(filename)

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := tmpl.Execute(file, data); err != nil {
		return err
	}

	return nil
}
