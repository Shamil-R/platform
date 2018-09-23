package schema

import (
	"gitlab/nefco/platform/codegen/template"
	"io/ioutil"
	"strings"

	"github.com/gobuffalo/packr"

	"github.com/vektah/gqlparser"
	"github.com/vektah/gqlparser/ast"
)

type Schema struct {
	*ast.Schema
}

func NewSchema(schema *ast.Schema) *Schema {
	return &Schema{schema}
}

func (s *Schema) ObjectTypes() []*ast.Definition {
	return s.filterTypes(ast.Object)
}

func (s *Schema) EnumTypes() []*ast.Definition {
	return s.filterTypes(ast.Enum)
}

func (s *Schema) InputTypes() []*ast.Definition {
	return s.ObjectTypes()
}

func (s *Schema) MutationTypes() []*ast.Definition {
	return s.ObjectTypes()
}

func (s *Schema) QueryTypes() []*ast.Definition {
	return s.ObjectTypes()
}

func (s *Schema) filterTypes(kind ast.DefinitionKind) []*ast.Definition {
	types := make([]*ast.Definition, 0)
	for _, def := range s.Schema.Types {
		if !strings.HasPrefix(def.Name, "__") && def.Kind == kind {
			types = append(types, def)
		}
	}
	return types
}

func Load(files ...string) (*Schema, error) {
	box := packr.NewBox("./graphql")

	source := make([]*ast.Source, 0, len(files)+1)
	source = append(source,
		&ast.Source{
			Name:  "directives",
			Input: box.String("directives.graphql"),
		},
	)

	for _, filename := range files {
		file, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}

		source = append(source,
			&ast.Source{
				Name:  filename,
				Input: string(file),
			},
		)
	}

	schema, gqlErr := gqlparser.LoadSchema(source...)
	if gqlErr != nil {
		return nil, gqlErr
	}

	return NewSchema(schema), nil
}

func Generate(cfg Config) error {
	box := packr.NewBox("./templates")

	tmpl, err := template.Read("schema", box)
	if err != nil {
		return err
	}

	schema, err := Load(cfg.Source)
	if err != nil {
		return err
	}

	if err := template.Execute(tmpl, schema, cfg.Generate); err != nil {
		return err
	}

	return nil
}
