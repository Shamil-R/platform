package build

type Code struct {
	PackageName string
	Imports     []*Import
	Schema      *Schema
}

type Import struct {
	Path  string
	Alias string
}
