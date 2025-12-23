package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

type lipglossWorkspaceModel struct {
	workspaces []string
	cursor     int
	selected   string
	action     string
	quitting   bool
	showBanner bool
}

func (m lipglossWorkspaceModel) Init() tea.Cmd {
	return nil
}

func (m lipglossWorkspaceModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			m.action = "quit"
			return m, tea.Quit
		case "c":
			m.action = "create"
			return m, tea.Quit
		case "d":
			if len(m.workspaces) > 0 {
				m.selected = m.workspaces[m.cursor]
				m.action = "delete"
				return m, tea.Quit
			}
		case "h":
			m.action = "dashboard"
			return m, tea.Quit
		case "g":
			if len(m.workspaces) > 0 {
				m.selected = m.workspaces[m.cursor]
				m.action = "goto"
				return m, tea.Quit
			}
		case "enter":
			if len(m.workspaces) > 0 {
				m.selected = m.workspaces[m.cursor]
				m.action = "select"
				return m, tea.Quit
			}
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.workspaces)-1 {
				m.cursor++
			}
		}
	}
	return m, cmd
}

func (m lipglossWorkspaceModel) View() string {
	if m.quitting {
		return ""
	}

	var s string
	if m.showBanner {

		bannerText := lipgloss.NewStyle().Foreground(lipgloss.Color("#E06C75")).Bold(true).Render(orbitBanner)
		centeredBanner := lipgloss.PlaceHorizontal(100, lipgloss.Center, bannerText)
		s += centeredBanner + "\n"

		subtitle := lipgloss.NewStyle().Foreground(lipgloss.Color("#5C6370")).Italic(true).Render("Keep your side projects in orbit 🚀")
		centeredSubtitle := lipgloss.PlaceHorizontal(100, lipgloss.Center, subtitle)
		s += centeredSubtitle + "\n\n"
	} else {
		s += "\n"
	}
	s += titleStyle.Render("Workspaces") + "\n\n"

	if len(m.workspaces) > 0 {
		s += m.renderTable() + "\n\n"
	} else {
		s += lipgloss.NewStyle().Foreground(mutedColor).Render("  No workspaces found. Press 'c' to create one.") + "\n\n"
	}

	helpBar := lipgloss.JoinHorizontal(lipgloss.Center,
		greenBtn.Render("c Create"),
		redBtn.Render("d Delete"),
		blueBtn.Render("g Goto"),
		blueBtn.Render("h Dashboard"),
		purpleBtn.Render("q Quit"),
	)
	s += helpBar + "\n"

	if len(m.workspaces) > 0 {
		s += lipgloss.NewStyle().Foreground(lipgloss.Color("#5C6370")).Render("\n  ↑/↓: Navigate  Enter: Select") + "\n"
	}

	return s
}

func (m lipglossWorkspaceModel) renderTable() string {

	const (
		workspaceWidth = 20
		lastModWidth   = 20
		pathWidth      = 55
	)


	var rows [][]string
	for _, w := range m.workspaces {
		name := truncateString(w[len(w)-len(w):], workspaceWidth)

		if idx := len(w); idx > 0 {

			for i := len(w) - 1; i >= 0; i-- {
				if w[i] == '/' {
					name = w[i+1:]
					break
				}
			}
		}

		lastMod := getLastModifiedTime(w)
		path := truncateString(w, pathWidth)

		rows = append(rows, []string{
			name,
			lastMod,
			path,
		})
	}


	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(primaryColor)).
		Headers("Workspace", "Last Modified", "Path").
		Rows(rows...).
		StyleFunc(func(row, col int) lipgloss.Style {
			var width int
			switch col {
			case 0:
				width = workspaceWidth
			case 1:
				width = lastModWidth
			case 2:
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

func selectWorkspaceLipgloss(workspaces []string, showBanner bool) (string, string) {
	m := lipglossWorkspaceModel{
		workspaces: workspaces,
		showBanner: showBanner,
		cursor:     0,
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	finalModel, _ := p.Run()
	result := finalModel.(lipglossWorkspaceModel)

	return result.selected, result.action
}
