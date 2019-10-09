package cmd

import (
	"fmt"
	"os"

	au "github.com/logrusorgru/aurora"

	"github.com/s4heid/goom/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewRootCmd returns a new cobra root command.
func NewRootCmd(configReader *config.Reader) *cobra.Command {
	var configPath string
	cobra.OnInitialize(func() {
		_ = configReader.SetConfig(configPath)
	})

	rootCmd := &cobra.Command{
		Use:     "goom",
		Short:   "A tool for opening urls from the command line.",
		Version: version,
	}

	var ioStreams = IOStreams{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "config file (default is $HOME/.goom.yaml)")
	rootCmd.AddCommand(
		NewOpenCmd(configReader, DefaultBrowser{}, ioStreams),
		NewShowCmd(configReader, ioStreams),
		NewVersionCmd(),
	)

	return rootCmd
}

// Execute runs the cobra root command.
func Execute() {
	configReader := &config.Reader{
		ViperConfig: viper.New(),
	}

	rootCmd := NewRootCmd(configReader)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(au.Red(fmt.Sprintf("execute command failed: %v", err)))
		os.Exit(1)
	}
}
