package gqlgen

import (
	"bytes"
	"gitlab/nefco/platform/codegen/helper"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Schema   helper.File
	Exec     helper.File
	Model    helper.File
	Resolver helper.File
	Dst      string
}

func Generate(cfg Config) error {
	var c struct {
		Schema string `yaml:"schema"`
		Exec   struct {
			Filename string `yaml:"filename"`
			Package  string `yaml:"package"`
		} `yaml:"exec"`
		Model struct {
			Filename string `yaml:"filename"`
			Package  string `yaml:"package"`
		} `yaml:"model"`
		Resolver struct {
			Filename string `yaml:"filename"`
			Type     string `yaml:"type"`
		} `yaml:"resolver"`
	}

	c.Schema = cfg.Schema.Path
	c.Exec.Filename = cfg.Exec.Path
	c.Exec.Package = cfg.Exec.Package()
	c.Model.Filename = cfg.Model.Path
	c.Model.Package = cfg.Model.Package()
	c.Resolver.Filename = cfg.Resolver.Path
	c.Resolver.Type = cfg.Resolver.Type

	b, err := yaml.Marshal(&c)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(b)

	if err := helper.WriteFile(cfg.Dst, buf); err != nil {
		return err
	}

	return nil
}
