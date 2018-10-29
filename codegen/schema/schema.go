package schema

import (
	"bytes"
	"io/ioutil"
	"strings"

	"github.com/gobuffalo/packr"
	"github.com/vektah/gqlparser"
	"github.com/vektah/gqlparser/ast"
)

type Schema struct {
	*ast.Schema
	types DefinitionList
}

func (s *Schema) Mutation() *Definition {
	if s.Schema.Mutation == nil {
		return nil
	}
	return &Definition{Definition: s.Schema.Mutation, schema: s}
}

func (s *Schema) Query() *Definition {
	if s.Schema.Query == nil {
		return nil
	}
	return &Definition{Definition: s.Schema.Query, schema: s}
}

func (s *Schema) Types() DefinitionList {
	if s.types != nil {
		return s.types
	}
	s.types = make(DefinitionList, 0, len(s.Schema.Types))
	for _, def := range s.Schema.Types {
		isInt := def.Name == "Int"
		isString := def.Name == "String"
		if !strings.HasPrefix(def.Name, "__") &&
			!isInt && !isString {
			s.types = append(s.types, &Definition{Definition: def, schema: s})
		}
	}
	return s.types
}

func LoadSchemaRaw(path string) (string, error) {
	box := packr.NewBox("./graphql")

	directivesRaw := box.Bytes("directives.graphql")

	buf := bytes.NewBuffer(directivesRaw)

	schemaRaw, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	if _, err := buf.Write(schemaRaw); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func ParseSchema(schemaRaw string) (*Schema, error) {
	source := &ast.Source{
		Name:  "schema",
		Input: schemaRaw,
	}

	s, err := gqlparser.LoadSchema(source)
	if err != nil {
		return nil, err
	}

	return &Schema{Schema: s}, nil
}

func LoadSchema(path string) (*Schema, error) {
	schemaRaw, err := LoadSchemaRaw(path)
	if err != nil {
		return nil, err
	}

	s, err := ParseSchema(schemaRaw)
	if err != nil {
		return nil, err
	}

	return s, nil
}
