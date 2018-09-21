package graphql

import (
	"github.com/gobuffalo/packr"
)

var Directives string

func init() {
	Directives = packr.NewBox("./").String("directives.graphql")
}
