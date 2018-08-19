package cmd

import (
	"fmt"
	"gitlab/nefco/platform/codegen"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(codegenCmd)
}

var codegenCmd = &cobra.Command{
	Use:   "codegen",
	Short: "codegen",
	Long:  `codegen`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := codegen.Generate(cfg.Codegen); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
