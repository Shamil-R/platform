package build

type Signature interface {
	Signature() string
}

type Item struct {
	Object
	Required bool
}

type List struct {
	Object
	Required bool
}

type Arg struct {
	Name string
}
