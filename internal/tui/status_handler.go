package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/henrynguci/orbit/internal/config"
)

func handleChangeStatus(cfg *config.Config, projectName string) {
	clearScreen()

	project, exists := cfg.Projects[projectName]
	if !exists {
		allProjects := config.GetAllProjects(cfg)
		for _, p := range allProjects {
			if p.Name == projectName {
				project = p
				exists = true
				break
			}
		}
		if !exists {
			printError("Project not found.")
			waitForEnter()
			return
		}
	}

	currentStatus := project.Status
	if currentStatus == "" {
		currentStatus = "not set"
	}

	fmt.Println(titleStyle.Render("Change Status"))
	fmt.Println()
	fmt.Printf("Project: %s\n", lipgloss.NewStyle().Foreground(secondaryColor).Bold(true).Render(projectName))
	fmt.Printf("Current status: %s\n", currentStatus)
	fmt.Println()

	var statusOptions []string
	if currentStatus != "active" {
		statusOptions = append(statusOptions, "active")
	}
	if currentStatus != "archived" {
		statusOptions = append(statusOptions, "archived")
	}
	if currentStatus != "done" {
		statusOptions = append(statusOptions, "done")
	}
	if currentStatus != "not set" {
		statusOptions = append(statusOptions, "not set")
	}

	newStatus, err := gumChoose(statusOptions, "Select new status:")
	if err != nil || newStatus == "" {
		return
	}

	statusValue := newStatus

	project.Name = projectName
	project.Status = statusValue
	cfg.Projects[projectName] = project
	config.Save(cfg)

	printSuccess(fmt.Sprintf("Status changed to '%s' for project '%s'", statusValue, projectName))
	waitForEnter()
}
