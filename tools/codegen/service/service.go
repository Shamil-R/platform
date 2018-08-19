package service

import "github.com/vektah/gqlparser/ast"

type Interface struct {
	*ast.Definition
}

func (d *Interface) Name() string {
	return d.Name()
}
