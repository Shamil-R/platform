package codegen

import (
	"gitlab/nefco/platform/codegen/schema"
	"gitlab/nefco/platform/codegen/service"
)

func GenerateSchema(cfg Config) error {
	if err := schema.Generate(cfg.Schema()); err != nil {
		return err
	}
	return nil
}

func GenerateService(cfg Config) error {
	if err := service.Generate(cfg.Service()); err != nil {
		return err
	}
	return nil
}
