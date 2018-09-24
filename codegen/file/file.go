package file

import (
	"io"
	"os"
	"path"
)

func Write(dst string, r io.Reader) error {
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
