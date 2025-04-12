package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ActionModel struct {
	CurrentScreen Screen
	grid          *Grid
	width         int
	height        int
}

func NewActionModel() *ActionModel {
	item := [][]GridItem{
		{
			ButtonItem{
				Title:   "Home button",
				Content: "Back to Home",
				ID:      "homeBtn",
				Style:   lipgloss.NewStyle(),
			},

			TextItem{
				Title:   "Wallet address",
				Content: "0x1234567890abcdef1234567890abcdef12345678",
				ID:      "walletAddress",
				Style:   lipgloss.NewStyle(),
			},
		},

		{
			ListItem{
				Title:    "Actions",
				Items:    []string{"Settings", "Apps", "Goal Commands", "Explore"},
				Cursor:   0,
				Selected: map[int]struct{}{},
				ID:       "actionsList",
				Style:    lipgloss.NewStyle(),
			},

			TextItem{
				Title:   "ActionPreview",
				Content: "This is a preview of the selected action.",
				ID:      "actionPreview",
				Style:   lipgloss.NewStyle(),
			},
		},
	}

	return &ActionModel{
		grid:          NewGrid(2, 2, item),
		CurrentScreen: ActionScreen,
	}
}

func (m ActionModel) Init() tea.Cmd {
	return nil
}

func (m ActionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "esc":
			if m.CurrentScreen == ActionScreen {
				return NewRootModel(), nil
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	// Maintain 4:1 aspect ratio
	if float64(m.width)/float64(m.height) > 4.0 {
		m.width = int(float64(m.height) * 4.0)
	} else {
		m.height = int(float64(m.width) / 4.0)
	}

	if m.width > 120 || m.height > 40 { // Max 80x24
		m.width = 120
		m.height = 40
	}

	return m, nil
}

func (m ActionModel) View() string {
	viewportWidth, viewportHeight := 80, 24 // Default size

	if m.width > 0 && m.height > 0 {
		viewportWidth, viewportHeight = m.width, m.height
	}

}
