package cmd

import (
	"github.com/henrynguci/orbit/internal/tui"
	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all projects in TUI",
	Run: func(cmd *cobra.Command, args []string) {
		tui.RunMainTUI()
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
