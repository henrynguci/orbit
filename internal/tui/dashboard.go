package tui

import (
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/henrynguci/orbit/internal/config"
)

type lipglossDashboardModel struct {
	workspaces []string
	projects   map[string]config.Project
	rawData    []dashboardRow
	cursor     int
	selected   string
	action     string
	quitting   bool
	menuOpen   bool
	menuIndex  int
}

func (m lipglossDashboardModel) Init() tea.Cmd {
	return nil
}

func (m lipglossDashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if m.menuOpen {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "up", "k":
				if m.menuIndex > 0 {
					m.menuIndex--
				}
			case "down", "j":
				if m.menuIndex < 2 {
					m.menuIndex++
				}
			case "enter":
				m.menuOpen = false
				actions := []string{"code", "cursor", "antigravityy"}
				m.action = "code_" + actions[m.menuIndex]
				return m, tea.Quit
			case "esc", "q":
				m.menuOpen = false
				return m, nil
			}
		}
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			m.action = "quit"
			return m, tea.Quit
		case "r", "esc":
			m.action = "return"
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.rawData)-1 {
				m.cursor++
			}
		case "g":
			if len(m.rawData) > 0 {
				row := m.rawData[m.cursor]
				m.selected = row.Path
				m.action = "goto"
				return m, tea.Quit
			}
		case "s":
			if len(m.rawData) > 0 {
				row := m.rawData[m.cursor]
				if row.Status != "none" {
					m.selected = row.Project
					m.action = "status"
					return m, tea.Quit
				}
			}
		case "m":
			if len(m.rawData) > 0 {
				row := m.rawData[m.cursor]
				if row.Status != "none" {
					m.selected = row.Path
					m.menuOpen = true
					m.menuIndex = 0
				}
			}
		case "enter":
			if len(m.rawData) > 0 {
				row := m.rawData[m.cursor]
				if row.Status != "none" {
					m.selected = row.Project
					m.action = "view"
					return m, tea.Quit
				}
			}
		}
	}
	return m, cmd
}

func (m lipglossDashboardModel) View() string {
	if m.quitting {
		return ""
	}

	var s string
	s += "\n"
	s += titleStyle.Render("Dashboard - All Projects") + "\n\n"

	if len(m.rawData) > 0 {
		s += m.renderTable() + "\n\n"
	} else {
		s += lipgloss.NewStyle().Foreground(mutedColor).Render("  No projects or workspaces found.") + "\n\n"
	}

	helpBar := lipgloss.JoinHorizontal(lipgloss.Center,
		blueBtn.Render("s Status"),
		blueBtn.Render("g Goto"),
		greenBtn.Render("m Code"),
		yellowBtn.Render("r Return"),
		purpleBtn.Render("q Quit"),
	)
	s += helpBar + "\n"

	if m.menuOpen {
		menuItems := []string{"open with code", "open with cursor", "open with antigravityy"}
		var menuLines []string
		for i, item := range menuItems {
			line := "  " + item
			if i == m.menuIndex {
				line = " > " + lipgloss.NewStyle().Foreground(primaryColor).Bold(true).Render(item)
			}
			menuLines = append(menuLines, line)
		}

		menuBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(1, 2).
			Render(lipgloss.JoinVertical(lipgloss.Left, menuLines...))

		s += "\n" + lipgloss.NewStyle().MarginLeft(4).Render(menuBox) + "\n"
	}

	if len(m.rawData) > 0 && !m.menuOpen {
		s += lipgloss.NewStyle().Foreground(lipgloss.Color("#5C6370")).Render("\n  ↑/↓: Navigate  Enter: View README") + "\n"
	}

	return s
}

func (m lipglossDashboardModel) renderTable() string {

	const (
		workspaceWidth = 12
		projectWidth   = 15
		statusWidth    = 10
		lastModWidth   = 18
		pathWidth      = 40
	)

	var rows [][]string
	for _, data := range m.rawData {

		statusText := data.Status
		if data.Status == "none" {
			statusText = "none"
		}

		projectText := data.Project

		lastMod := "none"
		if data.Project != "none" {
			lastMod = getLastModifiedTime(data.Path)
		}

		workspace := truncateString(data.Workspace, workspaceWidth)
		project := truncateString(projectText, projectWidth)
		path := truncateString(data.Path, pathWidth)

		rows = append(rows, []string{
			workspace,
			project,
			statusText,
			lastMod,
			path,
		})
	}

	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(primaryColor)).
		Headers("Workspace", "Project", "Status", "Last Modified", "Path").
		Rows(rows...).
		StyleFunc(func(row, col int) lipgloss.Style {
			var width int
			switch col {
			case 0:
				width = workspaceWidth
			case 1:
				width = projectWidth
			case 2:
				width = statusWidth
			case 3:
				width = lastModWidth
			case 4:
				width = pathWidth
			}

			style := lipgloss.NewStyle().Width(width).Padding(0, 1)

			if row == table.HeaderRow {
				return style.
					Bold(true).
					Foreground(primaryColor).
					BorderForeground(primaryColor).
					Align(lipgloss.Left)
			}

			var rowColor lipgloss.Color
			if row < len(m.rawData) {
				switch m.rawData[row].Status {
				case "active":
					rowColor = successColor
				case "archived":
					rowColor = warningColor
				case "done":
					rowColor = secondaryColor
				case "none":
					rowColor = mutedColor
				case "not set":
					rowColor = mutedColor
				default:
					rowColor = whiteColor
				}
			} else {
				rowColor = whiteColor
			}

			if row == m.cursor {
				return style.
					Background(primaryColor).
					Foreground(whiteColor).
					Bold(true)
			}

			return style.Foreground(rowColor)
		})

	return t.Render()
}

func selectDashboardLipgloss(workspaces []string, projects map[string]config.Project) (string, string) {
	_, rawData := prepareDashboardData(workspaces, projects)

	m := lipglossDashboardModel{
		workspaces: workspaces,
		projects:   projects,
		rawData:    rawData,
		cursor:     0,
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	finalModel, _ := p.Run()
	result := finalModel.(lipglossDashboardModel)

	return result.selected, result.action
}

func prepareDashboardData(workspaces []string, projects map[string]config.Project) ([][]string, []dashboardRow) {
	var rows [][]string
	var rawData []dashboardRow

	for _, w := range workspaces {
		wName := filepath.Base(w)
		wProjects := 0
		for _, p := range projects {
			if strings.HasPrefix(p.Path, w) {
				status := p.Status
				if status == "" {
					status = "not set"
				}
				lastMod := getLastModifiedTime(p.Path)

				rows = append(rows, []string{wName, p.Name, status, lastMod, p.Path})
				rawData = append(rawData, dashboardRow{Workspace: wName, Project: p.Name, Status: status, Path: p.Path})
				wProjects++
			}
		}
		if wProjects == 0 {
			rows = append(rows, []string{wName, "none", "", "none", w})
			rawData = append(rawData, dashboardRow{Workspace: wName, Project: "none", Status: "none", Path: w})
		}
	}

	return rows, rawData
}
