package main

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TableModel struct {
	table table.Model
}

type Column struct {
	Key   string
	Title string
	Width int
}

func NewTableModel(columns []table.Column, rows []table.Row) *TableModel {
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
		table.WithWidth(50),
	)

	t.SetStyles(defaultTableStyles())

	return &TableModel{table: t}
}

func defaultTableStyles() table.Styles {
	styles := table.DefaultStyles()

	styles.Cell = styles.Cell.
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 1)

	styles.Header = styles.Header.
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("205")).
		Bold(true)

	styles.Selected = styles.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(true)

	return styles
}

func (m *TableModel) Init() tea.Cmd {
	return nil
}

func (m *TableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.table = table.New(
			table.WithColumns(m.table.Columns()),
			table.WithRows(m.table.Rows()),
			table.WithFocused(true),
			table.WithHeight(msg.Height),
			table.WithWidth(msg.Width),
		)
		m.table.SetStyles(defaultTableStyles())
	}
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m *TableModel) View() string {
	return m.table.View()
}
