package cmd

import (
	"fmt"

	"github.com/henrynguci/orbit/internal/config"
	"github.com/henrynguci/orbit/internal/utils"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status [project]",
	Short: "Get project status",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]

		cfg, err := config.Load()
		if err != nil {
			utils.PrintError("Failed to load config: " + err.Error())
			return
		}

		project, exists := cfg.Projects[projectName]
		if !exists {
			projectPath := config.FindProjectPath(cfg, projectName)
			if projectPath == "" {
				utils.PrintError(fmt.Sprintf("Project '%s' not found", projectName))
				return
			}

			project = config.Project{
				Name:   projectName,
				Path:   projectPath,
				Status: "active",
			}
		}

		status := project.Status
		if status == "" {
			status = "active"
		}

		fmt.Printf("\n")
		fmt.Printf("  📁 Project: %s\n", projectName)
		if project.Alias != "" {
			fmt.Printf("  🏷️  Alias:   %s\n", project.Alias)
		}
		fmt.Printf("  📊 Status:  %s\n", status)
		fmt.Printf("  📍 Path:    %s\n", project.Path)
		fmt.Printf("\n")
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
