package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// ProjectModel is the first real accessible model for the user
// It will be the main entry point for the user to interact with the application
// selecting which screen to go to, from the main list of options

type ProjectModel struct {
	CurrentState SessionState
	options      []string
	cursor       int
	selected     map[int]struct{} // selected items
}

func NewProjectModel() *ProjectModel {
	return &ProjectModel{
		options: []string{
			"Settings",
			"Applications",
			"Cmd Goals",
			"Explore",
		},
		cursor:   0,
		selected: make(map[int]struct{}),
	}
}

func (m *ProjectModel) Init() tea.Cmd {
	return nil
}

func (m *ProjectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}
		case "enter", " ":

			// User can select only one option at a time
			// If user selects more than one option, the last one will be the selected one
			// and the previous ones will be unselected

			if _, ok := m.selected[m.cursor]; ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected = make(map[int]struct{})
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	return m, nil
}

// ProjectView has to show the preview of each option's view while the user is hovering it
// After the user selects an option, the focus has to move to the selected option's view

func (m *ProjectModel) View() string {
	// Header of the list
	s := "Which option do you want to select?\n\n"

	// Iterate over the list of options
	for i, option := range m.options {
		// Cursor position
		cursor := " " // No cursor
		if m.cursor == i {
			cursor = ">" // Cursor
		}

		selected := " "
		if _, ok := m.selected[i]; ok {
			selected = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, selected, option)
	}

	return s + "\nPress q to quit.\n"
}
