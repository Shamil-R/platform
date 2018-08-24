package build

import "strconv"

type Code struct {
	PackageName string
	Imports     []Import
}

type Import struct {
	Alias string
	Path  string
}

func (i *Import) Write() string {
	return i.Alias + " " + strconv.Quote(i.Path)
}
