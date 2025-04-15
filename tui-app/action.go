package main

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ActionModel struct {
	CurrentScreen Screen
	width         int
	height        int
	table         *TableModel
}

func NewActionModel() *ActionModel {
	columns := []table.Column{
		{Title: "Main Section", Width: 70},
		{Title: "Side Section", Width: 20},
	}

	rows := []table.Row{
		{"Header Table", "Address"},
		{"Main Table", "Side Table"},
	}

	tableModel := NewTableModel(columns, rows)

	return &ActionModel{
		CurrentScreen: ActionScreen,
		table:         tableModel,
	}
}

func (m *ActionModel) Init() tea.Cmd {
	return nil
}

func (m *ActionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

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
	}

	var cmd tea.Cmd
	updatedModel, cmd := m.table.Update(msg)
	if updatedTable, ok := updatedModel.(*TableModel); ok {
		m.table = updatedTable
	}
	return m, cmd
}

func (m *ActionModel) View() string {
	viewportWidth, viewportHeight := 120, 30 // Default size

	if m.width > 0 && m.height > 0 {
		viewportWidth, viewportHeight = m.width, m.height
	}

	switch m.CurrentScreen {
	case ActionScreen:
		tableView := m.table.View()
		container := createContainer(tableView, viewportWidth, viewportHeight)
		return lipgloss.NewStyle().Render(container)
	case HomeScreen:
		return NewRootModel().View()
	}

	return ""
}

func (m *ActionModel) Render() string {
	container := createContainer("Action Screen", m.width, m.height)
	return lipgloss.NewStyle().Render(container)
}
