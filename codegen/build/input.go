package build

import "github.com/vektah/gqlparser/ast"

type Input interface {
	Name() string
	Arg() string
}

type CreateInput struct {
	def *ast.Definition
}

func (i *CreateInput) Name() string {
	return i.def.Name + "CreateInput"
}

func (i *CreateInput) Arg() string {
	return "data: " + i.Name()
}

type UpdateInput struct {
	def *ast.Definition
}

func (i *UpdateInput) Name() string {
	return i.def.Name + "UpdateInput"
}

func (i *UpdateInput) Arg() string {
	return "data: " + i.Name()
}

type DeleteInput struct {
	def *ast.Definition
}

func (i *DeleteInput) Name() string {
	return i.def.Name + "DeleteInput"
}

func (i *DeleteInput) Arg() string {
	return "data: " + i.Name()
}

type WhereUniqueInput struct {
	def *ast.Definition
}

func (i *WhereUniqueInput) Name() string {
	return i.def.Name + "WhereUniqueInput"
}

func (i *WhereUniqueInput) Arg() string {
	return "where: " + i.Name()
}

type WhereInput struct {
	def *ast.Definition
}

func (i *WhereInput) Name() string {
	return i.def.Name + "WhereInput"
}

func (i *WhereInput) Arg() string {
	return "where: " + i.Name()
}
