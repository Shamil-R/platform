package schema

var DefaultConfig = Config{
	Source:   "schema.graphql",
	Generate: "schema_gen",
}

type Config struct {
	Source   string `mapstructure:"source"`
	Generate string `mapstructure:"generate"`
}
