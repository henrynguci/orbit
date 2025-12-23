package cmd

import (
	"fmt"

	"github.com/henrynguci/orbit/internal/config"
	"github.com/henrynguci/orbit/internal/utils"
	"github.com/spf13/cobra"
)

var aliasCmd = &cobra.Command{
	Use:   "alias [project] [alias]",
	Short: "Set an alias for a project",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		alias := args[1]

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

		project.Alias = alias
		cfg.Projects[projectName] = project
		cfg.Projects[alias] = project

		if err := config.Save(cfg); err != nil {
			utils.PrintError("Failed to save config: " + err.Error())
			return
		}

		utils.PrintSuccess(fmt.Sprintf("Alias '%s' set for project '%s'", alias, projectName))
	},
}

func init() {
	rootCmd.AddCommand(aliasCmd)
}
