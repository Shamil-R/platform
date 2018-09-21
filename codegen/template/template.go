package template

import (
	"fmt"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/gobuffalo/packr"
	"github.com/jinzhu/inflection"
)

func Read(box packr.Box) (*template.Template, error) {
	if len(box.List()) == 0 {
		return nil, fmt.Errorf("template: box empty in call to Parse")
	}
	funcs := sprig.TxtFuncMap()
	funcs["plural"] = inflection.Plural
	var tmpl *template.Template
	for _, name := range box.List() {
		var t *template.Template
		if tmpl == nil {
			tmpl = template.New(name).Funcs(funcs)
		}
		if name == tmpl.Name() {
			t = tmpl
		} else {
			t = tmpl.New(name).Funcs(funcs)
		}
		_, err := t.Parse(box.String(name))
		if err != nil {
			return nil, err
		}
	}

	return tmpl, nil
}
