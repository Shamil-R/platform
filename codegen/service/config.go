package service

var DefaultConfig = Config{
	Package:  "service",
	Filename: "service_gen.go",
}

type Config struct {
	Package  string
	Filename string
	Schema   string
	Includes string
}
