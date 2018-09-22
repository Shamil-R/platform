package build

import "gitlab/nefco/platform/codegen/schema"

type Code struct {
	PackageName string
	Imports     []*Import
	Schema      *schema.Schema
}

type Import struct {
	Path  string
	Alias string
}
