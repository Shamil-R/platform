package cmd

import (
	"fmt"
	"gitlab/nefco/platform/codegen"
	"os"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

var (
	genSchema   = false
	genGqlgen   = false
	genService  = false
	genServer   = false
	genClean    = false
	genInit		= false
)

func init() {
	codegenCmd.Flags().BoolVar(&genSchema, "schema", false, "generate schema")
	codegenCmd.Flags().BoolVar(&genGqlgen, "gqlgen", false, "generate gqlgen")
	codegenCmd.Flags().BoolVar(&genService, "service", false, "generate service")
	codegenCmd.Flags().BoolVar(&genServer, "server", false, "generate server")
	codegenCmd.Flags().BoolVar(&genClean, "clean", false, "generate clean")
	codegenCmd.Flags().BoolVar(&genInit, "init", false, "generate init")

	rootCmd.AddCommand(codegenCmd)
}

var codegenCmd = &cobra.Command{
	Use:   "codegen",
	Short: "codegen",
	Long:  `codegen`,
	Run: func(cmd *cobra.Command, args []string) {
		if genClean {
			if err := codegen.Clean(); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		if genSchema {
			if err := codegen.GenerateSchema(viper.GetViper()); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else if genGqlgen {
			if err := codegen.GenerateGqlgen(viper.GetViper()); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else if genService {
			if err := codegen.GenerateService(viper.GetViper()); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else if genInit {
			if err := codegen.GenerateInit(viper.GetViper()); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else  if genServer {
			if err := codegen.GenerateServer(viper.GetViper()); err != nil {
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
