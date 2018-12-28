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
		isBoolean := def.Name == "Boolean"
		if !strings.HasPrefix(def.Name, "__") &&
			!isInt && !isString && !isBoolean {
			s.types = append(s.types, &Definition{Definition: def, schema: s})
		}
	}
	return s.types
}

func writeDirectives(buf *bytes.Buffer) error {
	box := packr.NewBox("./graphql")

	directivesRaw := box.Bytes("directives.graphql")

	if _, err := buf.Write(directivesRaw); err != nil {
		return err
	}
	return nil
}

func LoadSchemaRaw(path string) (string, error) {
	buf := bytes.NewBufferString("")

	if err := writeDirectives(buf); err != nil {
		return "", err
	}

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

func ParseSchemaWithDirectives(schemaRaw string) (*Schema, error) {
	buf := bytes.NewBufferString("")

	if err := writeDirectives(buf); err != nil {
		return nil, err
	}

	if _, err := buf.WriteString(schemaRaw); err != nil {
		return nil, err
	}

	s, err := ParseSchema(buf.String())
	if err != nil {
		return nil, err
	}

	return s, nil
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
