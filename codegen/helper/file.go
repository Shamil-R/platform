package helper

import (
	"io"
	"os"
	"path"
)

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

func WriteFile(dst string, r io.Reader) error {
	dir := path.Dir(dst)

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	file, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := io.Copy(file, r); err != nil {
		return err
	}

	return nil
}
