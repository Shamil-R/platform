package service

func Generate(cfg Config) error {
	// tmpl, err := templates.Template()
	// if err != nil {
	// 	return err
	// }

	// schema, err := schema.Load(cfg.Schema)
	// if err != nil {
	// 	return err
	// }

	// code := &build.Code{
	// 	PackageName: cfg.Package(),
	// 	Imports: []*build.Import{
	// 		&build.Import{
	// 			Path:  "gitlab/nefco/platform/" + cfg.ModelPath(),
	// 			Alias: "model",
	// 		},
	// 	},
	// 	Schema: schema,
	// }

	// if err := template.Execute(tmpl, code, cfg.Filename); err != nil {
	// 	return err
	// }

	return nil
}
