package helper

import "path"

type File struct {
	Root string
	Path string
	Type string
}

func (f File) Dir() string {
	return path.Dir(f.Path)
}

func (f File) Filename() string {
	return path.Base(f.Path)
}

func (f File) Package() string {
	return path.Base(f.Dir())
}

func (f File) Import() string {
	return path.Join(f.Root, f.Dir())
}
