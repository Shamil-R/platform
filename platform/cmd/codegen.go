package cmd

import (
	"fmt"
	"gitlab/nefco/platform/codegen"
	"os"

	"github.com/spf13/cobra"
)

var (
	genSchema  = false
	genService = false
)

func init() {
	codegenCmd.Flags().BoolVar(&genSchema, "schema", false, "generate schema")
	codegenCmd.Flags().BoolVar(&genService, "service", false, "generate service")

	rootCmd.AddCommand(codegenCmd)
}

var codegenCmd = &cobra.Command{
	Use:   "codegen",
	Short: "codegen",
	Long:  `codegen`,
	Run: func(cmd *cobra.Command, args []string) {
		if genSchema {
			if err := codegen.GenerateSchema(cfg.Codegen); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		if genService {
			if err := codegen.GenerateService(cfg.Codegen); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	},
}
