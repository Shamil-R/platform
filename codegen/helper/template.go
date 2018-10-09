package helper

import (
	"fmt"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/gobuffalo/packr"
	"github.com/jinzhu/inflection"
)

func ReadTemplate(name string, box packr.Box) (*template.Template, error) {
	if len(box.List()) == 0 {
		return nil, fmt.Errorf("template: box empty in call to ReadTemplate")
	}

	funcs := sprig.TxtFuncMap()
	funcs["plural"] = inflection.Plural
	funcs["next"] = next

	tmpl := template.New(name + ".gotpl").Funcs(funcs)

	for _, n := range box.List() {
		var t *template.Template

		if n == tmpl.Name() {
			t = tmpl
		} else {
			t = tmpl.New(n)
		}

		if _, err := t.Parse(box.String(n)); err != nil {
			return nil, err
		}
	}

	return tmpl, nil
}

func next(i, l int) bool {
	return l-i > 1
}
