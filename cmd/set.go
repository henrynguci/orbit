package cmd

import (
	"fmt"
	"strings"

	"github.com/henrynguci/orbit/internal/config"
	"github.com/henrynguci/orbit/internal/utils"
	"github.com/spf13/cobra"
)

var validStatuses = []string{"active", "archived", "done", "not set"}

var setCmd = &cobra.Command{
	Use:   "set [project] [status]",
	Short: "Set project status",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		status := strings.ToLower(args[1])

		valid := false
		for _, s := range validStatuses {
			if status == s {
				valid = true
				break
			}
		}
		if !valid {
			utils.PrintError(fmt.Sprintf("Invalid status '%s'. Valid: %s", status, strings.Join(validStatuses, ", ")))
			return
		}

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
				Name: projectName,
				Path: projectPath,
			}
		}

		project.Status = status
		cfg.Projects[projectName] = project

		if err := config.Save(cfg); err != nil {
			utils.PrintError("Failed to save config: " + err.Error())
			return
		}

		utils.PrintSuccess(fmt.Sprintf("Project '%s' status set to %s", projectName, status))
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
}
