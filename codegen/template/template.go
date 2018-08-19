package template

import (
	"io/ioutil"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/jinzhu/inflection"
)

func Parse(text string) (*template.Template, error) {
	funcs := sprig.TxtFuncMap()
	funcs["plural"] = inflection.Plural

	tmpl, err := template.New("template").Funcs(funcs).Parse(text)
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

func Load(filename string) (*template.Template, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	tmpl, err := Parse(string(file))
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}
