package schema

import (
	"bytes"
	"io/ioutil"
	"strings"

	"github.com/gobuffalo/packr"
	"github.com/vektah/gqlparser"
	"github.com/vektah/gqlparser/ast"
)

const (
	ACTION_CREATE     = "create"
	ACTION_UPDATE     = "update"
	ACTION_DELETE     = "delete"
	ACTION_ITEM       = "item"
	ACTION_COLLECTION = "collection"
	ACTION_RELATION   = "relation"
)

type Schema struct {
	*ast.Schema
}

func (s *Schema) Mutation() *Definition {
	if s.Schema.Mutation == nil {
		return nil
	}
	return &Definition{s.Schema.Mutation, s}
}

func (s *Schema) Query() *Definition {
	if s.Schema.Query == nil {
		return nil
	}
	return &Definition{s.Schema.Query, s}
}

func (s *Schema) Types() DefinitionList {
	definitions := make(DefinitionList, 0, len(s.Schema.Types))
	for _, def := range s.Schema.Types {
		if !strings.HasPrefix(def.Name, "__") {
			definitions = append(definitions, &Definition{def, s})
		}
	}
	return definitions
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

	schema, err := gqlparser.LoadSchema(source)
	if err != nil {
		return nil, err
	}

	return &Schema{schema}, nil
}

func LoadSchema(path string) (*Schema, error) {
	schemaRaw, err := LoadSchemaRaw(path)
	if err != nil {
		return nil, err
	}

	schema, err := ParseSchema(schemaRaw)
	if err != nil {
		return nil, err
	}

	return schema, nil
}
