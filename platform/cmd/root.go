package cmd

import (
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	cfg     = DefaultConfig
)

var rootCmd = &cobra.Command{
	Use:   "platform",
	Short: "platform",
	Long:  `platform`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
		"config file (default is $HOME/.platform.yml)",
	)
}

func initConfig() {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		work, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.AddConfigPath(work)
		viper.SetConfigName(".platform")
	}

	viper.SetEnvPrefix("platform")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	for _, v := range configVars(cfg) {
		viper.BindEnv(v)
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		fmt.Println("Can't unmarshal config:", err)
		os.Exit(1)
	}
}
