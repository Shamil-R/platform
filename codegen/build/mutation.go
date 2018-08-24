package build

import "github.com/vektah/gqlparser/ast"

type Mutation interface {
	Name() string
}

type CreateMutation struct {
	def *ast.Definition
}

func (m *CreateMutation) Name() string {
	return "create" + m.def.Name
}

func (m *CreateMutation) Args() []*Arg {
	return nil
}

type UpdateMutation struct {
	def *ast.Definition
}

func (m *UpdateMutation) Name() string {
	return "update" + m.def.Name
}

type DeleteMutation struct {
	def *ast.Definition
}

func (m *DeleteMutation) Name() string {
	return "delete" + m.def.Name
}
