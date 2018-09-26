package code

type Code struct {
	PackageName string
	Imports     []*Import
}

func (c *Code) AddImport(path string, alias string) {
	c.Imports = append(c.Imports, &Import{path, alias})
}

type Import struct {
	Path  string
	Alias string
}
