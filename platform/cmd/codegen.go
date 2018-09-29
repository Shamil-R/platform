package cmd

import (
	"fmt"
	"gitlab/nefco/platform/codegen"
	"os"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

var (
	genGqlgen  = false
	genSchema  = false
	genService = false
)

func init() {
	codegenCmd.Flags().BoolVar(&genGqlgen, "gqlgen", false, "generate gqlgen")
	codegenCmd.Flags().BoolVar(&genSchema, "schema", false, "generate schema")
	codegenCmd.Flags().BoolVar(&genService, "service", false, "generate service")

	rootCmd.AddCommand(codegenCmd)
}

var codegenCmd = &cobra.Command{
	Use:   "codegen",
	Short: "codegen",
	Long:  `codegen`,
	Run: func(cmd *cobra.Command, args []string) {
		if genGqlgen {
			if err := codegen.GenerateGqlgen(viper.GetViper()); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else if genSchema {
			if err := codegen.GenerateSchema(viper.GetViper()); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else if genService {
			if err := codegen.GenerateService(viper.GetViper()); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
			if err := codegen.Generate(viper.GetViper()); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	},
}
