package build

import (
	"github.com/vektah/gqlparser/ast"
)

type Model struct {
	def *ast.Definition
}

func (m *Model) Name() string {
	return m.def.Name
}

func (m *Model) CreateInput() *CreateInput {
	return &CreateInput{m.def}
}

func (m *Model) UpdateInput() *UpdateInput {
	return &UpdateInput{m.def}
}

func (m *Model) DeleteInput() *DeleteInput {
	return &DeleteInput{m.def}
}

func (m *Model) WhereUniqueInput() *WhereUniqueInput {
	return &WhereUniqueInput{m.def}
}

func (m *Model) WhereInput() *WhereInput {
	return &WhereInput{m.def}
}

func (m *Model) CreateMutation() *CreateMutation {
	return &CreateMutation{m.def}
}

func (m *Model) UpdateMutation() *UpdateMutation {
	return &UpdateMutation{m.def}
}

func (m *Model) DeleteMutation() *DeleteMutation {
	return &DeleteMutation{m.def}
}
