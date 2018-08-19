package transform

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/jinzhu/inflection"

	"github.com/Masterminds/sprig"

	"github.com/vektah/gqlparser"
	"github.com/vektah/gqlparser/ast"
	"go.uber.org/zap"
)

func Transform(src, dest string) error {
	fmt.Println("transform")

	srcFile, err := ioutil.ReadFile(src)
	if err != nil {
		zap.L().Error("read file", zap.Error(err))
		return err
	}

	source := &ast.Source{
		Name:  "schema",
		Input: string(srcFile),
	}

	schema, e := gqlparser.LoadSchema(source)
	if e != nil {
		zap.L().Error("load schema", zap.Error(err))
		return err
	}

	types := make(map[string]*ast.Definition)

	for key, value := range schema.Types {
		if value.IsCompositeType() && !strings.HasPrefix(key, "__") {
			types[key] = value
		}
	}

	newSchema := ast.Schema{
		Types: types,
	}

	// tplFile, err := ioutil.ReadFile("../tools/transform/templates/schema.tpl")
	// if err != nil {
	// 	zap.L().Error("read file", zap.Error(err))
	// 	return err
	// }

	pattern := filepath.Join(
		"../tools/transform/templates",
		"*.tmpl",
	)

	funcs := sprig.TxtFuncMap()
	funcs["plural"] = inflection.Plural

	tmpl := template.Must(
		template.New("schema.tmpl").Funcs(funcs).ParseGlob(pattern),
	)
	if err != nil {
		zap.L().Error("template parse", zap.Error(err))
		return err
	}

	fmt.Println("schema")
	if err := tmpl.Execute(os.Stdout, newSchema); err != nil {
		zap.L().Error("execute", zap.Error(err))
		return err
	}

	return nil
}
