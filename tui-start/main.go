package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// CodeLab color palette
// 006883 - Teal Algo
// 31D8EE - Lightblue CodeLab
// 01DC94 - Green CodeLab
// 201F21 - Darkmode Algo version
// 575755 - Darkmode Vanilla

// Lipgloss styles
var (
	titleStyle = lipgloss.NewStyle().
			Width(40).
			Height(3).
			Padding(1).
			Margin(1, 0, 0, 0).
			Align(lipgloss.Center).
			Background(lipgloss.Color("#201F21")).
			Foreground(lipgloss.Color("#31D8EE")).
			Bold(true)

	subtitleStyle = lipgloss.NewStyle().
			Width(40).
			Height(2).
			Padding(1).
			Margin(0, 0, 1, 0).
			Align(lipgloss.Center).
			Background(lipgloss.Color("#201F21")).
			Foreground(lipgloss.Color("#01DC94"))

	messageStyle = lipgloss.NewStyle().
			Width(40).
			Height(2).
			Padding(1).
			Margin(1, 0, 0, 0).
			Align(lipgloss.Center).
			Background(lipgloss.Color("#201F21")).
			Foreground(lipgloss.Color("#575755"))
)

// Model
type model struct {
	screen   string
	cursor   int
	items    []string
	selected map[int]struct{}
}

// Init
func initialModel() model {
	return model{
		screen:   "home",
		cursor:   0,
		items:    []string{"settings", "applications", "base commands", "explore"},
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.screen == "home" {
				m.screen = "list"
			} else if m.screen == "list" {
				_, ok := m.selected[m.cursor]
				if ok {
					delete(m.selected, m.cursor)
				} else {
					m.selected[m.cursor] = struct{}{}
				}
			}
		case "backspace":
			if m.screen == "list" {
				m.screen = "home"
			}
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up":
			if m.screen == "list" && m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.screen == "list" && m.cursor < len(m.items)-1 {
				m.cursor++
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	switch m.screen {
	case "home":
		return homeView(m)
	case "list":
		return listView(m)
		// case "settings":
		// 	return settingsView(m)
		// case "applications":
		// 	return applicationsView(m)
		// case "base commands":
		// 	return baseCommandsView(m)
		// case "explore":
		// 	return exploreView(m)
	}

	return ""
}

func homeView(m model) string {
	title := titleStyle.Render("Welcome to the TUI!")
	subtitle := subtitleStyle.Render("Press 'Enter' to view the list")
	message := messageStyle.Render("Press 'Ctrl+C' or 'q' to quit")
	return fmt.Sprintf("%s\n%s\n%s", title, subtitle, message)
}

func listView(m model) string {
	title := titleStyle.Render("Select an option")
	subtitle := subtitleStyle.Render("Use 'arrows' to navigate and press 'Enter' to select an option")
	message := messageStyle.Render("Press 'Backspace' to go back to the home screen or 'Ctrl+C' to quit")
	s := fmt.Sprintf("%s\n\n", title)
	for i, item := range m.items {
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor
		}

		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked = "x" // selected
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, item)
	}

	return fmt.Sprintf("%s\n%s\n%s", s, subtitle, message)
}

// Main function
func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Println("Error starting TUI:", err)
	}
}
