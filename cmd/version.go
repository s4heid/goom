package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var version string

func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Display goom version information",
		Run: func(cmd *cobra.Command, args []string) {
			if version == "" {
				version = "0.0.0+dev"
			}
			fmt.Printf("goom/%s (%s; %s; %s/%s)",
				version,
				runtime.Version(),
				runtime.Compiler,
				runtime.GOOS,
				runtime.GOARCH,
			)
		},
	}
}

func init() {
	versionCmd := NewVersionCmd()
	rootCmd.AddCommand(versionCmd)
}
