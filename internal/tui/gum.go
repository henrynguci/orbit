package tui

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/henrynguci/orbit/internal/config"
)

var (
	primaryColor   = lipgloss.Color("#7C3AED")
	secondaryColor = lipgloss.Color("#00B9E8")
	successColor   = lipgloss.Color("#00FF00")
	warningColor   = lipgloss.Color("#E25822")
	errorColor     = lipgloss.Color("#FF0800")
	mutedColor     = lipgloss.Color("#6B7280")
	whiteColor     = lipgloss.Color("#FFFFFF")

	greenBtn = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#282C34")).
			Background(lipgloss.Color("#98C379")).
			Padding(0, 2).
			MarginRight(1).
			Bold(true)

	redBtn = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#282C34")).
		Background(lipgloss.Color("#E06C75")).
		Padding(0, 2).
		MarginRight(1).
		Bold(true)

	blueBtn = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#282C34")).
		Background(lipgloss.Color("#61AFEF")).
		Padding(0, 2).
		MarginRight(1).
		Bold(true)

	purpleBtn = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#282C34")).
			Background(lipgloss.Color("#C678DD")).
			Padding(0, 2).
			MarginRight(1).
			Bold(true)
	yellowBtn = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#282C34")).
			Background(lipgloss.Color("#F59E0B")).
			Padding(0, 2).
			MarginRight(1).
			Bold(true)

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(mutedColor)

	cellStyle = lipgloss.NewStyle().
			Foreground(whiteColor).
			Padding(0, 1)

	selectedStyle = lipgloss.NewStyle().
			Foreground(whiteColor).
			Background(primaryColor).
			Bold(true).
			Padding(0, 1)

	actionButtonStyle = lipgloss.NewStyle().
				Foreground(whiteColor).
				Background(lipgloss.Color("#374151")).
				Padding(0, 2).
				Margin(0, 1)

	activeButtonStyle = lipgloss.NewStyle().
				Foreground(whiteColor).
				Background(primaryColor).
				Bold(true).
				Padding(0, 2).
				Margin(0, 1)

	dangerButtonStyle = lipgloss.NewStyle().
				Foreground(whiteColor).
				Background(errorColor).
				Bold(true).
				Padding(0, 2).
				Margin(0, 1)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(whiteColor).
			Background(primaryColor).
			Padding(1, 4).
			Width(100).
			Align(lipgloss.Center).
			MarginBottom(1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(1, 2)
)

func gumChoose(items []string, header string) (string, error) {
	args := []string{"choose", "--header", header}
	args = append(args, items...)

	cmd := exec.Command("gum", args...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(out.String()), nil
}

func gumConfirm(prompt string) bool {
	cmd := exec.Command("gum", "confirm", prompt)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	return err == nil
}

func gumInput(placeholder string, header string) (string, error) {
	cmd := exec.Command("gum", "input", "--placeholder", placeholder, "--header", header)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(out.String()), nil
}

func printBanner() {
	banner := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF6B6B")).
		Bold(true).
		Render(orbitBanner)

	subtitle := subtitleStyle.Render("  Keep your side projects in orbit 🚀")

	fmt.Println(banner)
	fmt.Println(subtitle)
	fmt.Println()
}

func printSuccess(msg string) {
	style := lipgloss.NewStyle().Foreground(successColor).Bold(true)
	fmt.Println(style.Render(" " + msg))
}

func printError(msg string) {
	style := lipgloss.NewStyle().Foreground(errorColor).Bold(true)
	fmt.Println(style.Render("❌ " + msg))
}

func printInfo(msg string) {
	style := lipgloss.NewStyle().Foreground(secondaryColor)
	fmt.Println(style.Render("ℹ️  " + msg))
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func waitForReturnOrQuit() string {
	return handleReturnQuitKeyboard()
}

func RunMainTUI() {
	firstRun := true
	for {
		cfg, _ := config.Load()
		if cfg == nil {
			cfg = &config.Config{Workspaces: []string{}, Projects: make(map[string]config.Project)}
		}

		validWorkspaces := filterExistingWorkspaces(cfg.Workspaces)
		cfg.Workspaces = validWorkspaces
		config.Save(cfg)

		selected, action := selectWorkspaceLipgloss(validWorkspaces, firstRun)
		firstRun = false

		switch action {
		case "quit":
			clearScreen()
			return
		case "create":
			handleCreateWorkspace()
		case "delete":
			if selected != "" {
				handleDeleteWorkspace(cfg, selected)
			}
		case "dashboard":
			showDashboard()
		case "goto":
			if selected != "" {
				if err := os.Chdir(selected); err != nil {
					printError(fmt.Sprintf("Failed to access directory: %v", err))
					waitForReturnOrQuit()
					continue
				}
				shell := os.Getenv("SHELL")
				if shell == "" {
					shell = "/bin/bash"
				}
				clearScreen()
				syscall.Exec(shell, []string{shell}, os.Environ())
				os.Exit(0)
			}
		case "select":
			if selected != "" {
				handleWorkspaceViewWithTable(selected)
			}
		}
	}
}

func handleWorkspaceViewWithTable(workspace string) {
	for {
		cfg, _ := config.Load()
		projects := getProjectsInWorkspace(cfg, workspace)

		selected, action := selectProjectLipgloss(projects, workspace)

		switch {
		case action == "quit":
			clearScreen()
			os.Exit(0)
		case action == "return":
			return
		case action == "add":
			handleAddProjectToWorkspace(workspace)
		case action == "delete":
			if selected != "" {
				handleDeleteProject(cfg, selected)
			}
		case action == "status":
			if selected != "" {
				handleChangeStatus(cfg, selected)
			}
		case action == "view":
			if selected != "" {
				for _, p := range projects {
					if p.Name == selected {
						showProjectView(p)
						break
					}
				}
			}
		case action == "goto":
			if selected != "" {
				if err := os.Chdir(selected); err != nil {
					printError(fmt.Sprintf("Failed to access directory: %v", err))
					waitForReturnOrQuit()
					continue
				}
				shell := os.Getenv("SHELL")
				if shell == "" {
					shell = "/bin/bash"
				}
				clearScreen()
				syscall.Exec(shell, []string{shell}, os.Environ())
				os.Exit(0)
			}
		case strings.HasPrefix(action, "code_"):
			tool := strings.TrimPrefix(action, "code_")
			if selected != "" {
				cmd := exec.Command(tool, selected)
				cmd.Run()
			}
		}
	}
}

func printWorkspaceTable(workspaces []string) {
	colWidths := []int{25, 20, 50}

	header := lipgloss.JoinHorizontal(lipgloss.Left,
		headerStyle.Width(colWidths[0]).Render("Workspace"),
		headerStyle.Width(colWidths[1]).Render("Last Modified"),
		headerStyle.Width(colWidths[2]).Render("Path"),
	)
	fmt.Println(header)

	divider := lipgloss.NewStyle().Foreground(mutedColor).Render(strings.Repeat("─", colWidths[0]+colWidths[1]+colWidths[2]))
	fmt.Println(divider)

	for _, w := range workspaces {
		name := filepath.Base(w)
		lastMod := getLastModified(w)

		row := lipgloss.JoinHorizontal(lipgloss.Left,
			cellStyle.Width(colWidths[0]).Render(truncate(name, colWidths[0]-2)),
			cellStyle.Width(colWidths[1]).Render(lastMod),
			cellStyle.Width(colWidths[2]).Render(truncate(w, colWidths[2]-2)),
		)
		fmt.Println(row)
	}
}

func printProjectTable(projects []config.Project) {
	colWidths := []int{20, 15, 20, 40}

	header := lipgloss.JoinHorizontal(lipgloss.Left,
		headerStyle.Width(colWidths[0]).Render("Project"),
		headerStyle.Width(colWidths[1]).Render("Status"),
		headerStyle.Width(colWidths[2]).Render("Last Modified"),
		headerStyle.Width(colWidths[3]).Render("Path"),
	)
	fmt.Println(header)

	divider := lipgloss.NewStyle().Foreground(mutedColor).Render(strings.Repeat("─", colWidths[0]+colWidths[1]+colWidths[2]+colWidths[3]))
	fmt.Println(divider)

	for _, p := range projects {
		status := p.Status
		if status == "" {
			status = "not set"
		}
		lastMod := getLastModified(p.Path)

		statusStyle := lipgloss.NewStyle()
		switch status {
		case "active":
			statusStyle = statusStyle.Foreground(successColor)
		case "archived":
			statusStyle = statusStyle.Foreground(warningColor)
		case "done":
			statusStyle = statusStyle.Foreground(secondaryColor)
		case "not set":
			statusStyle = statusStyle.Foreground(mutedColor)
		}

		row := lipgloss.JoinHorizontal(lipgloss.Left,
			cellStyle.Width(colWidths[0]).Render(truncate(p.Name, colWidths[0]-2)),
			statusStyle.Inherit(cellStyle).Width(colWidths[1]).Render(status),
			cellStyle.Width(colWidths[2]).Render(lastMod),
			cellStyle.Width(colWidths[3]).Render(truncate(p.Path, colWidths[3]-2)),
		)
		fmt.Println(row)
	}
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func getLastModified(path string) string {
	info, err := os.Stat(path)
	if err != nil {
		return "Unknown"
	}

	modTime := info.ModTime()
	now := time.Now()

	if modTime.Year() == now.Year() && modTime.YearDay() == now.YearDay() {
		return fmt.Sprintf("Today %02d:%02d", modTime.Hour(), modTime.Minute())
	}

	if modTime.Year() == now.Year() && modTime.YearDay() == now.YearDay()-1 {
		return fmt.Sprintf("Yesterday %02d:%02d", modTime.Hour(), modTime.Minute())
	}

	return modTime.Format("02/01/2006 15:04")
}

func handleCreateWorkspace() {
	clearScreen()

	fmt.Println(titleStyle.Render("Create Workspace"))
	fmt.Println()

	path, err := gumInput("~/workspace", "Enter workspace path:")
	if err != nil {
		return
	}

	if path == "" {
		path = "~/workspace/orbit_ws1"
	}

	if strings.HasPrefix(path, "~") {
		home, _ := os.UserHomeDir()
		path = filepath.Join(home, path[1:])
	}
	absPath, _ := filepath.Abs(path)

	cfg, _ := config.Load()
	if cfg == nil {
		cfg = &config.Config{Workspaces: []string{}, Projects: make(map[string]config.Project)}
	}

	newName := filepath.Base(absPath)
	for _, w := range cfg.Workspaces {
		if filepath.Base(w) == newName {
			printError(fmt.Sprintf("Workspace with name '%s' already exists.", newName))
			waitForEnter()
			return
		}
	}

	if err := os.MkdirAll(absPath, 0755); err != nil {
		printError(fmt.Sprintf("Failed to create workspace: %v", err))
		waitForEnter()
		return
	}

	cfg.Workspaces = appendUnique(cfg.Workspaces, absPath)
	config.Save(cfg)

	printSuccess(fmt.Sprintf("Workspace created at %s", absPath))

	if gumConfirm("Create a project in this workspace?") {
		handleAddProjectToWorkspace(absPath)
	}
}

func handleDeleteWorkspace(cfg *config.Config, workspacePath string) {
	clearScreen()

	fmt.Println(titleStyle.Render("Delete Workspace"))
	fmt.Println()
	fmt.Println(lipgloss.NewStyle().Foreground(warningColor).Bold(true).Render("⚠️  Warning: This action cannot be undone!"))
	fmt.Println()
	fmt.Println("Please enter the exact path to confirm deletion:")
	fmt.Println(lipgloss.NewStyle().Foreground(secondaryColor).Bold(true).Render(workspacePath))
	fmt.Println()

	confirmPath, err := gumInput("", "Enter workspace path:")
	if err != nil || confirmPath == "" {
		return
	}

	if confirmPath != workspacePath {
		printError("Path does not match. Deletion cancelled.")
		waitForEnter()
		return
	}

	if !gumConfirm(fmt.Sprintf("Delete workspace '%s'?", filepath.Base(workspacePath))) {
		printInfo("Deletion cancelled.")
		waitForEnter()
		return
	}

	var newWorkspaces []string
	for _, w := range cfg.Workspaces {
		if w != workspacePath {
			newWorkspaces = append(newWorkspaces, w)
		}
	}
	cfg.Workspaces = newWorkspaces

	var projectsToRemove []string
	for name, p := range cfg.Projects {
		if strings.HasPrefix(p.Path, workspacePath) {
			projectsToRemove = append(projectsToRemove, name)
		}
	}
	for _, name := range projectsToRemove {
		delete(cfg.Projects, name)
	}

	config.Save(cfg)

	if gumConfirm("Also delete workspace files from disk?") {
		os.RemoveAll(workspacePath)
		printSuccess("Workspace and files deleted.")
	} else {
		printSuccess("Workspace removed from Orbit (files kept on disk).")
	}
	waitForEnter()
}

func showDashboard() {
	for {
		cfg, _ := config.Load()
		if cfg == nil {
			printInfo("No projects found.")
			waitForReturnOrQuit()
			return
		}

		// Get all projects with their status
		allProjects := config.GetAllProjects(cfg)
		projectsMap := make(map[string]config.Project)
		for _, p := range allProjects {
			projectsMap[p.Name] = p
		}

		selected, action := selectDashboardLipgloss(cfg.Workspaces, projectsMap)

		switch {
		case action == "quit":
			clearScreen()
			os.Exit(0)
		case action == "return":
			return
		case action == "goto":
			if selected != "" {
				if err := os.Chdir(selected); err != nil {
					printError(fmt.Sprintf("Failed to access directory: %v", err))
					waitForReturnOrQuit()
					continue
				}
				shell := os.Getenv("SHELL")
				if shell == "" {
					shell = "/bin/bash"
				}
				clearScreen()
				syscall.Exec(shell, []string{shell}, os.Environ())
				os.Exit(0)
			}
		case strings.HasPrefix(action, "code_"):
			tool := strings.TrimPrefix(action, "code_")
			if selected != "" {
				cmd := exec.Command(tool, selected)
				if tool == "antigravityy" {
				}
				cmd.Run()
			}
		case action == "status":
			if selected != "" {
				handleChangeStatus(cfg, selected)
			}
		case action == "view":
			if selected != "" {
				if p, ok := cfg.Projects[selected]; ok {
					showProjectView(p)
				}
			}
		}
	}
}

func getProjectsInWorkspace(cfg *config.Config, workspace string) []config.Project {
	var projects []config.Project
	if cfg == nil {
		return projects
	}

	seenPaths := make(map[string]bool)

	for _, p := range cfg.Projects {
		if strings.HasPrefix(p.Path, workspace) && !seenPaths[p.Path] {
			projects = append(projects, p)
			seenPaths[p.Path] = true
		}
	}

	entries, err := os.ReadDir(workspace)
	if err != nil {
		return projects
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		projectPath := filepath.Join(workspace, entry.Name())
		if seenPaths[projectPath] {
			continue
		}

		status := "not set"
		if existing, exists := cfg.Projects[entry.Name()]; exists {
			status = existing.Status
		}

		projects = append(projects, config.Project{
			Name:   entry.Name(),
			Path:   projectPath,
			Status: status,
		})
		seenPaths[projectPath] = true
	}

	return projects
}

func handleSelectProject(cfg *config.Config, projects []config.Project, workspace string) {
	if len(projects) == 0 {
		return
	}

	items := make([]string, len(projects))
	for i, p := range projects {
		status := p.Status
		if status == "" {
			status = "active"
		}
		items[i] = fmt.Sprintf("%s %s │ %s", getStatusIcon(status), p.Name, p.Path)
	}

	choice, err := gumChoose(items, "Select project to view:")
	if err != nil || choice == "" {
		return
	}

	for _, p := range projects {
		if strings.Contains(choice, p.Name) {
			showProjectView(p)
			return
		}
	}
}

func showProjectView(project config.Project) {
	showProjectReadme(project.Name, project.Path)
}

func handleAddProjectToWorkspace(workspace string) {
	clearScreen()

	fmt.Println(titleStyle.Render("Add Project"))
	fmt.Println(subtitleStyle.Render(fmt.Sprintf("Workspace: %s", workspace)))
	fmt.Println()

	cfg, _ := config.Load()
	if cfg == nil {
		cfg = &config.Config{Workspaces: []string{}, Projects: make(map[string]config.Project)}
	}

	projectName, err := gumInput("my-project", "Enter project name:")
	if err != nil || projectName == "" {
		return
	}

	if _, exists := cfg.Projects[projectName]; exists {
		printError(fmt.Sprintf("Project '%s' already exists.", projectName))
		waitForEnter()
		return
	}

	projectPath := filepath.Join(workspace, "project", projectName)
	if _, err := os.Stat(projectPath); err == nil {
		printError(fmt.Sprintf("Directory '%s' already exists.", projectPath))
		waitForEnter()
		return
	}

	cloneRepo := gumConfirm("Clone a repository?")

	if cloneRepo {
		cloneURL, err := gumInput("https://github.com/user/repo.git", "Enter clone URL:")
		if err != nil || cloneURL == "" {
			return
		}

		repoPath := filepath.Join(projectPath, "repo")
		os.MkdirAll(filepath.Join(projectPath, "docs"), 0755)
		os.MkdirAll(filepath.Join(projectPath, "secret"), 0755)

		fmt.Println()
		printInfo("Cloning repository...")

		err = cloneRepository(cloneURL, repoPath)
		if err != nil {
			printError(fmt.Sprintf("Clone failed: %v", err))
			waitForEnter()
			return
		}

		saveProject(workspace, projectName, projectPath)
		printSuccess(fmt.Sprintf("Project '%s' created with cloned repo", projectName))
	} else {
		dirs := []string{
			filepath.Join(projectPath, "repo"),
			filepath.Join(projectPath, "docs"),
			filepath.Join(projectPath, "secret"),
		}
		for _, dir := range dirs {
			os.MkdirAll(dir, 0755)
		}

		saveProject(workspace, projectName, projectPath)
		printSuccess(fmt.Sprintf("Project '%s' created at %s", projectName, projectPath))
	}
	waitForEnter()
}

func handleDeleteProject(cfg *config.Config, projectName string) {
	clearScreen()

	project, exists := cfg.Projects[projectName]
	if !exists {
		printError("Project not found.")
		waitForEnter()
		return
	}

	fmt.Println(titleStyle.Render("Delete Project"))
	fmt.Println()
	fmt.Println(lipgloss.NewStyle().Foreground(warningColor).Bold(true).Render("⚠️  Warning: This action cannot be undone!"))
	fmt.Println()
	fmt.Println("Please enter the exact path to confirm deletion:")
	fmt.Println(lipgloss.NewStyle().Foreground(secondaryColor).Bold(true).Render(project.Path))
	fmt.Println()

	confirmPath, err := gumInput("", "Enter project path:")
	if err != nil || confirmPath == "" {
		return
	}

	if confirmPath != project.Path {
		printError("Path does not match. Deletion cancelled.")
		waitForEnter()
		return
	}

	if !gumConfirm(fmt.Sprintf("Delete project '%s'?", projectName)) {
		printInfo("Deletion cancelled.")
		waitForEnter()
		return
	}

	delete(cfg.Projects, projectName)
	config.Save(cfg)

	if gumConfirm("Also delete project files from disk?") {
		os.RemoveAll(project.Path)
		printSuccess("Project and files deleted.")
	} else {
		printSuccess("Project removed from Orbit (files kept on disk).")
	}
	waitForEnter()
}


func showProjectReadme(name, path string) {
	readmePath := filepath.Join(path, "repo", "README.md")
	if _, err := os.Stat(readmePath); os.IsNotExist(err) {
		readmePath = filepath.Join(path, "README.md")
		if _, err := os.Stat(readmePath); os.IsNotExist(err) {
			alternatives := []string{"readme.md", "Readme.md", "README.MD"}
			found := false
			for _, alt := range alternatives {
				altPath := filepath.Join(path, "repo", alt)
				if _, e := os.Stat(altPath); e == nil {
					readmePath = altPath
					found = true
					break
				}
				altPath = filepath.Join(path, alt)
				if _, e := os.Stat(altPath); e == nil {
					readmePath = altPath
					found = true
					break
				}
			}
			if !found {
				printInfo("README.md not found")
				return
			}
		}
	}

	cmd := exec.Command("glow", "-p", readmePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func saveProject(workspacePath, projectName, projectPath string) {
	cfg, _ := config.Load()
	if cfg == nil {
		cfg = &config.Config{Workspaces: []string{}, Projects: make(map[string]config.Project)}
	}

	cfg.Workspaces = appendUnique(cfg.Workspaces, workspacePath)
	cfg.Projects[projectName] = config.Project{
		Name:   projectName,
		Path:   projectPath,
		Status: "active",
	}
	config.Save(cfg)
}

func appendUnique(slice []string, item string) []string {
	for _, s := range slice {
		if s == item {
			return slice
		}
	}
	return append(slice, item)
}

func filterExistingWorkspaces(workspaces []string) []string {
	var valid []string
	for _, w := range workspaces {
		if _, err := os.Stat(w); err == nil {
			valid = append(valid, w)
		}
	}
	return valid
}

func getStatusIcon(status string) string {
	return ""
}

type handleReturnQuitKeyboardModel struct {
	action string
}

func (m handleReturnQuitKeyboardModel) Init() tea.Cmd { return nil }
func (m handleReturnQuitKeyboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			m.action = "return"
			return m, tea.Quit
		case "q":
			m.action = "quit"
			return m, tea.Quit
		case "ctrl+c":
			os.Exit(0)
		}
	}
	return m, nil
}
func (m handleReturnQuitKeyboardModel) View() string {
	return "\n" + lipgloss.JoinHorizontal(lipgloss.Center,
		yellowBtn.Render("r Return"),
		purpleBtn.Render("q Quit"),
	) + "\n"
}

func handleReturnQuitKeyboard() string {
	p := tea.NewProgram(handleReturnQuitKeyboardModel{})
	m, _ := p.Run()
	res := m.(handleReturnQuitKeyboardModel)
	if res.action == "quit" {
		clearScreen()
		os.Exit(0)
	}
	return "return"
}

func waitForEnter() {
	waitForReturnOrQuit()
}
