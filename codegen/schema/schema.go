//go:generate gorunpkg github.com/gobuffalo/packr/packr

package schema

import (
	"fmt"
	"gitlab/nefco/platform/codegen/build"
	"gitlab/nefco/platform/codegen/template"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/gobuffalo/packr"
	"github.com/jinzhu/inflection"
	"github.com/vektah/gqlparser"

	"github.com/vektah/gqlparser/ast"
)

func Generate(cfg Config) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	src := path.Join(wd, cfg.Source)
	dest := path.Join(wd, cfg.Generate)

	if err := os.MkdirAll(dest, os.ModePerm); err != nil {
		return err
	}

	box := packr.NewBox("./templates")

	tmpl, err := template.Parse(box.String("model.tpl"))
	if err != nil {
		return err
	}

	file, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	source := &ast.Source{
		Name:  "schema",
		Input: string(file),
	}

	gqlSchema, gqlErr := gqlparser.LoadSchema(source)
	if gqlErr != nil {
		return gqlErr
	}

	schema := build.NewSchema(gqlSchema)

	for _, def := range schema.Types() {
		file, err := os.Create(path.Join(dest, def.Name+".graphql"))
		if err != nil {
			return err
		}

		err = tmpl.Execute(file, def)

		file.Close()

		if err != nil {
			return err
		}
	}

	f, err := ioutil.ReadFile(path.Join(dest, "User.graphql"))
	if err != nil {
		return err
	}

	sourceUser := &ast.Source{
		Name:  "User",
		Input: string(f),
	}

	sources := make([]*ast.Source, 0, 2)
	sources = append(sources, source, sourceUser)

	gqlSchema, gqlErr = gqlparser.LoadSchema(sources...)
	if gqlErr != nil {
		return gqlErr
	}

	for _, f := range gqlSchema.Mutation.Fields {
		fmt.Println("!", f.Name, f.Type)
		for _, a := range f.Arguments {
			fmt.Println("#", a.Name, a.Type)
		}
	}

	return nil
}

func NewSchema(schema *ast.Schema) *Schema {
	return &Schema{
		types: toDefinitions(schema.Types),
	}
}

func Load(filename string) (*Schema, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	source := &ast.Source{
		Name:  "schema",
		Input: string(file),
	}

	schema, gqlerr := gqlparser.LoadSchema(source)
	if gqlerr != nil {
		return nil, gqlerr
	}

	return NewSchema(schema), nil
}

func Transform(src, dest string) error {
	schema, err := Load(src)
	if err != nil {
		return err
	}

	box := packr.NewBox("./templates")

	tmpl, err := template.Parse(box.String("schema.tpl"))
	if err != nil {
		return err
	}

	file, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := tmpl.Execute(file, schema); err != nil {
		return err
	}

	return nil
}

func toDefinitions(list map[string]*ast.Definition) []*Definition {
	definitions := make([]*Definition, 0, len(list))
	for _, def := range list {
		if def.IsCompositeType() && !strings.HasPrefix(def.Name, "__") {
			definition := &Definition{
				Definition: def,
				fields:     toFields(def.Fields),
				Input: Input{
					Create:      &CreateInput{def},
					Update:      &UpdateInput{def},
					WhereUnique: &WhereUniqueInput{def},
					Where:       &WhereInput{def},
				},
				Mutation: Mutation{
					Create: "create" + def.Name,
					Update: "update" + def.Name,
					Delete: "delete" + def.Name,
				},
				Query: Query{
					Item: strings.ToLower(def.Name),
					List: inflection.Plural(strings.ToLower(def.Name)),
				},
			}
			definitions = append(definitions, definition)
		}
	}
	return definitions
}

func toFields(list ast.FieldList) []*Field {
	fields := make([]*Field, len(list))
	for i, def := range list {
		fields[i] = &Field{
			FieldDefinition: def,
		}
	}
	return fields
}
