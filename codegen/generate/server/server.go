package server

import (
	"gitlab/nefco/platform/codegen/helper"
)

type Config struct {
	Schema  helper.File
	Server  helper.File
	Exec    helper.File
	Service helper.File
}

func Generate(cfg Config) error {
	// box := packr.NewBox("./templates")

	// tmpl, err := helper.ReadTemplate("server", box)
	// if err != nil {
	// 	return err
	// }

	// s, err := schema.LoadSchema(cfg.Schema.Path)
	// if err != nil {
	// 	return err
	// }

	// data := &struct {
	// 	*Config
	// 	Types schema.DefinitionList
	// }{
	// 	Config: &cfg,
	// 	Types:  s.Types().Objects(),
	// }

	// buff := &bytes.Buffer{}

	// if err := tmpl.Execute(buff, data); err != nil {
	// 	return err
	// }

	// if err := helper.WriteFile(cfg.Server.Path, buff); err != nil {
	// 	return err
	// }

	return nil
}
