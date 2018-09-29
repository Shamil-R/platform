package gqlgen

import (
	"bytes"
	"gitlab/nefco/platform/codegen/file"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Schema   string         `yaml:"schema"`
	Exec     ConfigExec     `yaml:"exec"`
	Model    ConfigModel    `yaml:"model"`
	Resolver ConfigResolver `yaml:"resolver"`
	Output   string         `yaml:"-"`
}

type ConfigExec struct {
	Filename string `yaml:"filename"`
	Package  string `yaml:"package"`
}

type ConfigModel struct {
	Filename string `yaml:"filename"`
	Package  string `yaml:"package"`
}

type ConfigResolver struct {
	Filename string `yaml:"filename"`
	Type     string `yaml:"type"`
}

func Generate(cfg Config) error {
	b, err := yaml.Marshal(&cfg)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(b)
	if err := file.Write(cfg.Output, buf); err != nil {
		return err
	}
	return nil
}
