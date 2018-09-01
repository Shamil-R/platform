package service

import "path"

type Config struct {
	Filename      string
	Schema        string
	ModelFilename string
}

func (c Config) Package() string {
	return path.Base(path.Dir(c.Filename))
}

func (c Config) ModelPath() string {
	return path.Dir(c.ModelFilename)
}
