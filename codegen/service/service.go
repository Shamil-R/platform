package service

import (
	"bytes"
	"fmt"
	"gitlab/nefco/platform/codegen/file"
	"gitlab/nefco/platform/codegen/schema"
	"gitlab/nefco/platform/codegen/service/code"
	"gitlab/nefco/platform/codegen/template"
	"path"
	"strings"

	"github.com/huandu/xstrings"

	"github.com/jinzhu/inflection"

	"github.com/gobuffalo/packr"
)

type Config struct {
	SchemaPath     string
	ServiceDir     string
	ServicePackage string
	ModelImport    string
}

type Code struct {
	*code.Code
	Schema *schema.Schema
}

type ServiceGenerator func() error

func Generate(cfg Config) error {
	box := packr.NewBox("./templates")

	tmpl, err := template.Read("service", box)
	if err != nil {
		return err
	}

	schema, err := schema.Load(cfg.SchemaPath)
	if err != nil {
		return err
	}

	code := &Code{
		Code: &code.Code{
			PackageName: cfg.ServicePackage,
		},
		Schema: schema,
	}
	code.AddImport("context", "context")
	code.AddImport(cfg.ModelImport, "model")

	buff := &bytes.Buffer{}

	if err := tmpl.Execute(buff, code); err != nil {
		return err
	}

	filename := path.Join(cfg.ServiceDir + "service_gen.go")

	if err := file.Write(filename, buff); err != nil {
		return err
	}

	// mssqlCfg := mssql.Config{
	// 	ServiceDir:     cfg.ServiceDir,
	// 	ServicePackage: cfg.ServicePackage,
	// 	ModelImport:    cfg.ModelImport,
	// }

	if err := generateObjectService(cfg, schema); err != nil {
		return err
	}

	return nil
}

type Action struct {
	Name  string
	Field *schema.FieldDefinition
}

func generateObjectService(cfg Config, sch *schema.Schema) error {
	box := packr.NewBox("./templates")

	tmpl, err := template.Read("service_struct", box)
	if err != nil {
		return err
	}

	for _, def := range sch.Types().ForMutation() {
		lname := xstrings.FirstRuneToLower(def.Name)
		names := map[string]string{
			"create":     fmt.Sprintf("create%s", def.Name),
			"update":     fmt.Sprintf("update%s", def.Name),
			"delete":     fmt.Sprintf("delete%s", def.Name),
			"item":       lname,
			"collection": inflection.Plural(lname),
		}

		actions := make([]*Action, 0, len(names))

		for key, value := range names {
			action := &Action{
				Name:  key,
				Field: sch.MutationAndQueryFields().ByName(value),
			}
			actions = append(actions, action)
		}

		data := &struct {
			*code.Code
			TypeName string
			Actions  []*Action
		}{
			Code:     code.New(cfg.ServicePackage),
			TypeName: def.Name,
			Actions:  actions,
		}
		data.AddImport("context", "context")
		data.AddImport(cfg.ModelImport, "model")

		buff := &bytes.Buffer{}

		if err := tmpl.Execute(buff, data); err != nil {
			return err
		}

		serviceName := strings.ToLower(def.Name) + "_service_gen.go"

		filename := path.Join(cfg.ServiceDir, serviceName)

		if err := file.Write(filename, buff); err != nil {
			return err
		}
	}

	return nil
}
