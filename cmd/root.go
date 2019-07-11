package cmd

import (
	"fmt"
	"os"

	. "github.com/logrusorgru/aurora"
	. "github.com/s4heid/goom/config"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var (
	configPath   = ""
	ConfigReader = Reader{ViperConfig: viper.New()}
	rootCmd      = &cobra.Command{
		Use:     "goom",
		Short:   "A tool for opening urls from the command line based on an alias matching the url.",
		Long:    "A tool for opening urls from the command line based on an alias matching the url.",
		Version: "0.0.1",
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(Red(fmt.Sprintf("execute command failed: %v", err)))
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initializeConfig)
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "config file (default is $HOME/.goom.yaml)")
}

func initializeConfig() {
	if ConfigReader.ViperConfig.ConfigFileUsed() == "" {
		err := ConfigReader.InitConfig(configPath)
		if err != nil {
			fmt.Println(Red(err))
			os.Exit(1)
		}
	}
}