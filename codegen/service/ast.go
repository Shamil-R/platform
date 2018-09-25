package service

import (
	"gitlab/nefco/platform/codegen/schema"
)

type Service struct {
	PackageName string
	Imports     []*Import
	Schema      *schema.Schema
}

func NewService(cfg Config, schema *schema.Schema) *Service {
	return &Service{
		PackageName: cfg.Package(),
		Imports: []*Import{
			&Import{
				Path:  "context",
				Alias: "context",
			},
			&Import{
				Path:  "gitlab/nefco/platform/" + cfg.ModelPath(),
				Alias: "model",
			},
		},
		Schema: schema,
	}
}

type Import struct {
	Path  string
	Alias string
}
