package helper

import (
	"fmt"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/gobuffalo/packr"
	"github.com/huandu/xstrings"
	"github.com/jinzhu/inflection"
)

func ReadTemplate(name string, boxes ...packr.Box) (*template.Template, error) {
	m := make(map[string]string)
	for _, box := range boxes {
		for _, n := range box.List() {
			m[n] = box.String(n)
		}
	}

	if len(m) == 0 {
		return nil, fmt.Errorf("template: boxes empty in call to ReadTemplate")
	}

	funcs := sprig.TxtFuncMap()
	funcs["plural"] = inflection.Plural
	funcs["firstRuneToUpper"] = xstrings.FirstRuneToUpper
	funcs["firstRuneToLower"] = xstrings.FirstRuneToLower

	tmpl := template.New(name + ".gotpl").Funcs(funcs)

	for n, text := range m {
		var t *template.Template

		if n == tmpl.Name() {
			t = tmpl
		} else {
			t = tmpl.New(n)
		}

		if _, err := t.Parse(text); err != nil {
			return nil, err
		}
	}

	return tmpl, nil
}
