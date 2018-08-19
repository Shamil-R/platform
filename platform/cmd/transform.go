package cmd

import (
	"fmt"
	"gitlab/nefco/platform/tools/graphql"

	"github.com/spf13/cobra"
)

var (
	source      string
	destination string
)

func init() {
	transformCmd.Flags().StringVar(&source, "source", "", "source")
	transformCmd.MarkFlagRequired("source")
	transformCmd.Flags().StringVar(&destination, "destination", "", "destination")
	transformCmd.MarkFlagRequired("destination")

	rootCmd.AddCommand(transformCmd)
}

var transformCmd = &cobra.Command{
	Use:   "transform",
	Short: "transform",
	Long:  `transform`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("s=%s, d=%s\n", source, destination)
		if err := graphql.Transform(source, destination); err != nil {
			fmt.Println("error: ", err)
		}
	},
}
