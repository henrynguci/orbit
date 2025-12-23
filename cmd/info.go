package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/henrynguci/orbit/internal/config"
	"github.com/henrynguci/orbit/internal/utils"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info [project]",
	Short: "Show project README.md",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]

		cfg, err := config.Load()
		if err != nil {
			utils.PrintError("Failed to load config: " + err.Error())
			return
		}

		projectPath := config.FindProjectPath(cfg, projectName)
		if projectPath == "" {
			utils.PrintError(fmt.Sprintf("Project '%s' not found", projectName))
			return
		}

		readmePath := filepath.Join(projectPath, "repo", "README.md")
		if _, err := os.Stat(readmePath); os.IsNotExist(err) {
			readmePath = filepath.Join(projectPath, "README.md")
			if _, err := os.Stat(readmePath); os.IsNotExist(err) {
				alternatives := []string{"readme.md", "Readme.md", "README.MD"}
				found := false
				for _, alt := range alternatives {
					altPath := filepath.Join(projectPath, "repo", alt)
					if _, err := os.Stat(altPath); err == nil {
						readmePath = altPath
						found = true
						break
					}
					altPath = filepath.Join(projectPath, alt)
					if _, err := os.Stat(altPath); err == nil {
						readmePath = altPath
						found = true
						break
					}
				}
				if !found {
					utils.PrintError(fmt.Sprintf("README.md not found in %s", projectPath))
					return
				}
			}
		}

		glowCmd := exec.Command("glow", "-p", readmePath)
		glowCmd.Stdin = os.Stdin
		glowCmd.Stdout = os.Stdout
		glowCmd.Stderr = os.Stderr
		glowCmd.Run()
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
