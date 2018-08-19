package cmd

import (
	"gitlab/nefco/platform/tools/codegen"

	"github.com/fatih/structs"
)

var DefaultConfig = Config{
	Codegen: codegen.DefaultConfig,
}

type Config struct {
	Codegen codegen.Config `mapstructure:"codegen"`
}

func configVars(s interface{}) []string {
	res := make([]string, 0, 1)

	fields := structs.Fields(s)

	for _, field := range fields {
		tag := field.Tag("mapstructure")

		if structs.IsStruct(field.Value()) {
			arr := configVars(field.Value())

			for _, t := range arr {
				res = append(res, tag+"."+t)
			}
		} else {
			res = append(res, tag)
		}
	}

	return res
}
