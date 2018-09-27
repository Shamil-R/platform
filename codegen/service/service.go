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

type ServiceGenerator func() error

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
	ActionName string
	FieldName  string
	Field      *schema.FieldDefinition
}

func actions(name string, fields schema.FieldList) []*Action {
	return []*Action{
		&Action{
			ActionName: "create",
			FieldName:  fmt.Sprintf("create%s", name),
			Field:      fields.ByName(fmt.Sprintf("create%s", name)),
		},
		&Action{
			ActionName: "update",
			FieldName:  fmt.Sprintf("update%s", name),
			Field:      fields.ByName(fmt.Sprintf("update%s", name)),
		},
		&Action{
			ActionName: "delete",
			FieldName:  fmt.Sprintf("delete%s", name),
			Field:      fields.ByName(fmt.Sprintf("delete%s", name)),
		},
		&Action{
			ActionName: "item",
			FieldName:  xstrings.FirstRuneToLower(name),
			Field:      fields.ByName(xstrings.FirstRuneToLower(name)),
		},
		&Action{
			ActionName: "collection",
			FieldName:  inflection.Plural(xstrings.FirstRuneToLower(name)),
			Field:      fields.ByName(inflection.Plural(xstrings.FirstRuneToLower(name))),
		},
	}
}

type Service struct {
	Name    string
	Actions []*Action
}

func generateServiceInterface(cfg Config, sch *schema.Schema) error {
	box := packr.NewBox("./templates")

	tmpl, err := template.Read("service_interface", box)
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

	tmpl, err := template.Read("service_struct", box)
	if err != nil {
		return err
	}

	for _, def := range sch.Types().ForMutation() {
		data := &struct {
			*code.Code
			TypeName string
			Actions  []*Action
		}{
			Code:     code.New(cfg.ServicePackage),
			TypeName: def.Name,
			Actions:  actions(def.Name, sch.MutationAndQueryFields()),
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
