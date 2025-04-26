package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	zone "github.com/lrstanley/bubblezone"
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
				if m.items[m.cursor] == "settings" {
					m.screen = "settings"
				} else if m.items[m.cursor] == "applications" {
					m.screen = "applications"
				} else if m.items[m.cursor] == "base commands" {
					m.screen = "commands"
				} else if m.items[m.cursor] == "explore" {
					m.screen = "explore"
				}
				_, ok := m.selected[m.cursor]
				if ok {
					delete(m.selected, m.cursor)
				} else {
					m.selected[m.cursor] = struct{}{}
				}
			}
		case "esc":
			if m.screen == "home" {
				return m, tea.Quit
			} else if m.screen == "list" {
				m.screen = "home"
			} else if m.screen == "settings" {
				m.screen = "list"
			} else if m.screen == "applications" {
				m.screen = "list"
			} else if m.screen == "commands" {
				m.screen = "list"
			} else if m.screen == "explore" {
				m.screen = "list"
			}
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.items)-1 {
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
	case "settings":
		return settingsView()
	case "applications":
		return applicationsView()
	case "commands":
		return baseCommandsView()
	case "explore":
		return exploreView()
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
	subtitle := subtitleStyle.Render("Use 'Arrows' to navigate and press 'Enter' to select an option")
	message := messageStyle.Render("Press 'Esc' to go back to the home screen or 'Ctrl+C' to quit")
	s := fmt.Sprintf("%s\n\n", title)
	for i, item := range m.items {
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor
			item = lipgloss.NewStyle().Foreground(lipgloss.Color("#01DC94")).Render(item)
		}

		_, ok := m.selected[i]
		if ok {
			m.Update(item)
		}

		s += fmt.Sprintf("%s %s %s\n", cursor, "-", item)
	}

	return fmt.Sprintf("%s\n%s\n%s", s, subtitle, message)
}

func settingsView() string {
	title := titleStyle.Render("Settings")

	return title
}

func applicationsView() string {
	title := titleStyle.Render("Applications")

	return title
}

func baseCommandsView() string {
	title := titleStyle.Render("Base Commands")

	return title
}

func exploreView() string {
	title := titleStyle.Render("Explore")

	return title
}

// Main function
func main() {
	m := initialModel()
	zone.NewGlobal()
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseAllMotion())
	if err := p.Start(); err != nil {
		fmt.Println("Error starting TUI:", err)
	}
}
