package model

var DefaultConfig = Config{
	Package:  "model",
	Filename: "model_get.go",
}

type Config struct {
	Package  string `mapstructure:"package"`
	Filename string `mapstructure:"filename"`
}
