package service

import (
	"bytes"
	"fmt"
	"gitlab/nefco/platform/codegen/file"
	"gitlab/nefco/platform/codegen/schema"
	"gitlab/nefco/platform/codegen/service/code"
	"gitlab/nefco/platform/codegen/service/mssql"
	codegentemplate "gitlab/nefco/platform/codegen/template"
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

type ActionGenerate func(action string, field *schema.FieldDefinition) (string, error)

func Generate(cfg Config) error {
	schema, err := schema.Load(cfg.SchemaPath)
	if err != nil {
		return err
	}

	if err := generateServiceInterface(cfg, schema); err != nil {
		return err
	}

	if err := generateServiceStruct(cfg, schema); err != nil {
		return err
	}

	return nil
}

type Action struct {
	Action string
	Field  *schema.FieldDefinition
}

func actions(name string, fields schema.FieldList) []*Action {
	return []*Action{
		&Action{
			Action: "create",
			Field:  fields.ByName(fmt.Sprintf("create%s", name)),
		},
		&Action{
			Action: "update",
			Field:  fields.ByName(fmt.Sprintf("update%s", name)),
		},
		&Action{
			Action: "delete",
			Field:  fields.ByName(fmt.Sprintf("delete%s", name)),
		},
		&Action{
			Action: "item",
			Field:  fields.ByName(xstrings.FirstRuneToLower(name)),
		},
		&Action{
			Action: "collection",
			Field:  fields.ByName(inflection.Plural(xstrings.FirstRuneToLower(name))),
		},
	}
}

type Service struct {
	Name    string
	Actions []*Action
}

func generateServiceInterface(cfg Config, sch *schema.Schema) error {
	box := packr.NewBox("./templates")

	tmpl, err := codegentemplate.Read("service_interface", box)
	if err != nil {
		return err
	}

	types := sch.Types().ForMutation()

	services := make([]*Service, len(types))

	for i, def := range types {
		service := &Service{
			Name:    def.Name,
			Actions: actions(def.Name, sch.MutationAndQueryFields()),
		}
		services[i] = service
	}

	data := &struct {
		*code.Code
		Services []*Service
	}{
		Code:     code.New(cfg.ServicePackage),
		Services: services,
	}
	data.AddImport("context", "context")
	data.AddImport(cfg.ModelImport, "model")

	buff := &bytes.Buffer{}

	if err := tmpl.Execute(buff, data); err != nil {
		return err
	}

	filename := path.Join(cfg.ServiceDir + "service_gen.go")

	if err := file.Write(filename, buff); err != nil {
		return err
	}

	return nil
}

func generateServiceStruct(cfg Config, sch *schema.Schema) error {
	box := packr.NewBox("./templates")

	var actionGenerate ActionGenerate = mssql.Generate

	tmpl, err := codegentemplate.Read("service_struct", box)
	if err != nil {
		return err
	}

	for _, def := range sch.Types().ForMutation() {
		data := &struct {
			*code.Code
			TypeName       string
			Actions        []*Action
			ActionGenerate ActionGenerate
		}{
			Code:           code.New(cfg.ServicePackage),
			TypeName:       def.Name,
			Actions:        actions(def.Name, sch.MutationAndQueryFields()),
			ActionGenerate: actionGenerate,
		}
		data.AddImport("context", "context")
		data.AddImport(cfg.ModelImport, "model")
		data.AddImport("github.com/jmoiron/sqlx", "sqlx")

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
