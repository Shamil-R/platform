package stub

import (
	"github.com/gobuffalo/packr"
)

type stub struct{}

func New() *stub {
	return &stub{}
}

func (s *stub) Name() string {
	return "stub"
}

func (s *stub) Box() packr.Box {
	return packr.NewBox("./templates")
}
