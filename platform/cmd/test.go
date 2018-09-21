package cmd

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/spf13/cobra"
	"github.com/vektah/gqlparser"
	"github.com/vektah/gqlparser/ast"
)

func init() {
	rootCmd.AddCommand(testCmd)
}

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "test",
	Long:  `test`,
	Run: func(cmd *cobra.Command, args []string) {
		fileA, err := ioutil.ReadFile("./a.graphql")
		if err != nil {
			log.Fatal(err)
		}

		sourceA := &ast.Source{
			Name:  "a",
			Input: string(fileA),
		}

		fileB, err := ioutil.ReadFile("./b.graphql")
		if err != nil {
			log.Fatal(err)
		}

		sourceB := &ast.Source{
			Name:  "b",
			Input: string(fileB),
		}

		schema, gqlErr := gqlparser.LoadSchema(sourceA, sourceB)
		if gqlErr != nil {
			log.Fatal(gqlErr)
		}

		for k, v := range schema.Types {
			fmt.Println(k, v.Name)
		}
	},
}
