package cmd

import (
	"gitlab/nefco/platform/app"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run",
	Long:  `run`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := app.Run(viper.GetViper()); err != nil {
			panic(err)
		}
	},
}
