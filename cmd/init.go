package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/henrynguci/orbit/internal/config"
	"github.com/henrynguci/orbit/internal/utils"
	"github.com/spf13/cobra"
)

var projectName string

var initCmd = &cobra.Command{
	Use:   "init [path]",
	Short: "Initialize a new orbit workspace",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]

		if path[0] == '~' {
			home, err := os.UserHomeDir()
			if err != nil {
				utils.PrintError("Failed to get home directory: " + err.Error())
				return
			}
			path = filepath.Join(home, path[1:])
		}

		absPath, err := filepath.Abs(path)
		if err != nil {
			utils.PrintError("Failed to get absolute path: " + err.Error())
			return
		}

		cfg, err := config.Load()
		if err != nil {
			cfg = &config.Config{
				Workspaces: []string{},
				Projects:   make(map[string]config.Project),
			}
		}

		if projectName != "" {
			projectPath := filepath.Join(absPath, "project", projectName)
			dirs := []string{
				filepath.Join(projectPath, "repo"),
				filepath.Join(projectPath, "docs"),
				filepath.Join(projectPath, "secret"),
			}

			for _, dir := range dirs {
				if err := os.MkdirAll(dir, 0755); err != nil {
					utils.PrintError(fmt.Sprintf("Failed to create directory %s: %s", dir, err.Error()))
					return
				}
			}

			cfg.Projects[projectName] = config.Project{
				Name:   projectName,
				Path:   projectPath,
				Status: "active",
			}

			utils.PrintSuccess(fmt.Sprintf("Project '%s' created at %s", projectName, projectPath))
		} else {
			if err := os.MkdirAll(absPath, 0755); err != nil {
				utils.PrintError(fmt.Sprintf("Failed to create workspace: %s", err.Error()))
				return
			}

			utils.PrintSuccess(fmt.Sprintf("Workspace initialized at %s", absPath))
		}

		workspaceExists := false
		for _, w := range cfg.Workspaces {
			if w == absPath {
				workspaceExists = true
				break
			}
		}
		if !workspaceExists {
			cfg.Workspaces = append(cfg.Workspaces, absPath)
		}

		if err := config.Save(cfg); err != nil {
			utils.PrintError("Failed to save config: " + err.Error())
			return
		}
	},
}

func init() {
	initCmd.Flags().StringVarP(&projectName, "project", "p", "", "Project name")
	rootCmd.AddCommand(initCmd)
}
