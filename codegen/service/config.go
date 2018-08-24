package service

var DefaultConfig = Config{
	Package:  "service",
	Filename: "service.gen.go",
}

type Config struct {
	Package  string `mapstructure:"package"`
	Filename string `mapstructure:"filename"`
}
