package graphql

import (
	"io/ioutil"

	"github.com/vektah/gqlparser"

	"github.com/vektah/gqlparser/ast"
)

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
