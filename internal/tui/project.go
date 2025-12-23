package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/henrynguci/orbit/internal/config"
)

type lipglossProjectModel struct {
	projects  []config.Project
	workspace string
	cursor    int
	selected  string
	action    string
	quitting  bool
	menuOpen  bool
	menuIndex int
}

func (m lipglossProjectModel) Init() tea.Cmd {
	return nil
}

func (m lipglossProjectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		case "a":
			m.action = "add"
			return m, tea.Quit
		case "d":
			if len(m.projects) > 0 {
				m.selected = m.projects[m.cursor].Name
				m.action = "delete"
				return m, tea.Quit
			}
		case "s":
			if len(m.projects) > 0 {
				m.selected = m.projects[m.cursor].Name
				m.action = "status"
				return m, tea.Quit
			}
		case "g":
			if len(m.projects) > 0 {
				m.selected = m.projects[m.cursor].Path
				m.action = "goto"
				return m, tea.Quit
			}
		case "m":
			if len(m.projects) > 0 {
				m.selected = m.projects[m.cursor].Path
				m.menuOpen = true
				m.menuIndex = 0
			}
		case "r", "esc":
			m.action = "return"
			return m, tea.Quit
		case "enter":
			if len(m.projects) > 0 {
				m.selected = m.projects[m.cursor].Name
				m.action = "view"
				return m, tea.Quit
			}
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.projects)-1 {
				m.cursor++
			}
		}
	}
	return m, cmd
}

func (m lipglossProjectModel) View() string {
	if m.quitting {
		return ""
	}

	var s strings.Builder
	s.WriteString("\n")
	s.WriteString(titleStyle.Render("Workspace: "+m.workspace) + "\n")
	s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#5C6370")).Italic(true).Render(m.workspace) + "\n\n")

	if len(m.projects) > 0 {
		s.WriteString(m.renderTable() + "\n\n")
	} else {
		s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#5C6370")).Render("  No projects found. Press 'a' to add one.") + "\n\n")
	}

	helpBar := lipgloss.JoinHorizontal(lipgloss.Center,
		greenBtn.Render("a Add"),
		redBtn.Render("d Delete"),
		blueBtn.Render("s Status"),
		blueBtn.Render("g Goto"),
		greenBtn.Render("m Code"),
		yellowBtn.Render("r Return"),
		purpleBtn.Render("q Quit"),
	)
	s.WriteString(helpBar + "\n")

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

		s.WriteString("\n" + lipgloss.NewStyle().MarginLeft(4).Render(menuBox) + "\n")
	}

	if len(m.projects) > 0 && !m.menuOpen {
		s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#5C6370")).Render("\n  ↑/↓: Navigate  Enter: View README") + "\n")
	}

	return s.String()
}

func (m lipglossProjectModel) renderTable() string {

	const (
		projectWidth = 15
		statusWidth  = 15
		lastModWidth = 20
		pathWidth    = 45
	)


	var rows [][]string
	for _, p := range m.projects {
		status := p.Status
		if status == "" {
			status = "active"
		}

		var statusText string
		switch status {
		case "active":
			statusText = lipgloss.NewStyle().Foreground(successColor).Render(status)
		case "archived":
			statusText = lipgloss.NewStyle().Foreground(warningColor).Render(status)
		case "done":
			statusText = lipgloss.NewStyle().Foreground(secondaryColor).Render(status)
		default:
			statusText = status
		}

		lastMod := getLastModifiedTime(p.Path)

		project := truncateString(p.Name, projectWidth)
		path := truncateString(p.Path, pathWidth)

		rows = append(rows, []string{
			project,
			statusText,
			lastMod,
			path,
		})
	}


	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(primaryColor)).
		Headers("Project", "Status", "Last Modified", "Path").
		Rows(rows...).
		StyleFunc(func(row, col int) lipgloss.Style {
			var width int
			switch col {
			case 0:
				width = projectWidth
			case 1:
				width = statusWidth
			case 2:
				width = lastModWidth
			case 3:
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


			if row == m.cursor {
				return style.
					Background(primaryColor).
					Foreground(whiteColor).
					Bold(true)
			}


			return style.Foreground(whiteColor)
		})

	return t.Render()
}

func selectProjectLipgloss(projects []config.Project, workspace string) (string, string) {
	m := lipglossProjectModel{
		projects:  projects,
		workspace: workspace,
		cursor:    0,
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	finalModel, _ := p.Run()
	result := finalModel.(lipglossProjectModel)

	return result.selected, result.action
}
