package schema

type Config struct {
	Source   string `mapstructure:"source"`
	Generate string `mapstructure:"generate"`
}
