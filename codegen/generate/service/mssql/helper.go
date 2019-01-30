package mssql

import (
	"fmt"
	"gitlab/nefco/platform/codegen/generate/service/mssql/query"
)

func logQuery(query query.Query) {
	fmt.Println(query.Query())
	fmt.Println(query.Arg())
}
