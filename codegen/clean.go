package codegen

import (
	"fmt"
	"io/ioutil"
	"os"
)

var directory = "app/"
var deprecated map[string]bool = map[string]bool{"app.go": true}

func Clean() error {
	dir, err := ioutil.ReadDir(directory)
	if err != nil {
		return err
	}
	for _, d := range dir {
		if deprecated[d.Name()] {
			continue
		}
		err = os.RemoveAll(fmt.Sprintf("%s%s", directory, d.Name()))
		if err != nil {
			return err
		}
	}
	return nil
}
