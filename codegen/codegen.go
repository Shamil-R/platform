package codegen

import (
	"gitlab/nefco/platform/codegen/schema"
)

func Generate(cfg Config) error {
	if err := schema.Generate(cfg.Schema()); err != nil {
		return err
	}

	return nil
}
