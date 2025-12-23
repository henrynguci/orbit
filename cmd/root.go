package cmd

import (
	"fmt"
	"os"

	"github.com/henrynguci/orbit/internal/tui"
	"github.com/spf13/cobra"
)

var version = "1.0.0"

var rootCmd = &cobra.Command{
	Use:     "orbit",
	Short:   "Keep your side projects in orbit 🚀",
	Version: version,
	Run: func(cmd *cobra.Command, args []string) {
		tui.RunMainTUI()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
