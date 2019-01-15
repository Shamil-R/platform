package schema

import (
	"github.com/vektah/gqlparser/ast"
)

type Arguments interface {
	Arguments() ArgumentList
}

type Argument struct {
	*ast.Argument
	value *Value
}

func (a *Argument) Value() *Value {
	if a.value == nil {
		a.value = &Value{Value: a.Argument.Value}
	}
	return a.value
}

type ArgumentList []*Argument

type argumentListFilter func(arg *Argument) bool

func (l ArgumentList) size() int {
	return len(l)
}

func (l ArgumentList) filter(filter argumentListFilter) ArgumentList {
	args := make(ArgumentList, 0, len(l))
	for _, arg := range l {
		if filter(arg) {
			args = append(args, arg)
		}
	}
	return args
}

func (l ArgumentList) first(filter argumentListFilter) *Argument {
	r := l.filter(filter)
	if r.size() == 0 {
		return nil
	}
	return r[0]
}

func (l ArgumentList) IsNotEmpty() bool {
	return l.size() > 0
}

func (l ArgumentList) ByName(name string) *Argument {
	fn := func(arg *Argument) bool {
		return arg.Name == name
	}
	return l.first(fn)
}

func (l ArgumentList) Data() *Argument {
	return l.ByName("data")
}

func (l ArgumentList) Where() *Argument {
	return l.ByName("where")
}
