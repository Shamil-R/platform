package codegen

import (
	"gitlab/nefco/platform/codegen/schema"
	"os"
	"path"
)

func Generate(cfg Config) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	cfg.Schema.Source = path.Join(wd, cfg.Schema.Source)
	cfg.Schema.Transform = path.Join(wd, cfg.Schema.Transform)

	err = schema.Transform(cfg.Schema.Source, cfg.Schema.Transform)
	if err != nil {
		return err
	}

	return nil
}
