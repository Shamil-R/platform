//go:generate gorunpkg github.com/gobuffalo/packr/packr

package schema

import (
	"gitlab/nefco/platform/codegen/template"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gobuffalo/packr"
	"github.com/jinzhu/inflection"
	"github.com/vektah/gqlparser"

	"github.com/vektah/gqlparser/ast"
)

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
