package template

import (
	"io/ioutil"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/jinzhu/inflection"
)

func Load(filename string) (*template.Template, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	funcs := sprig.TxtFuncMap()
	funcs["plural"] = inflection.Plural

	tmpl, err := template.New("tpl").Funcs(funcs).Parse(string(file))
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}
