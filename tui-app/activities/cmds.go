package main

import (
	tea "github.com/charmbracelet/bubbletea"

	zone "github.com/lrstanley/bubblezone"
)

type CmdsModel struct {
	zone zone.Manager
	id   string
}

func NewCmdsModel() *CmdsModel {
	return &CmdsModel{
		id: zone.NewPrefix(),
	}
}

func (m *CmdsModel) Init() tea.Cmd {
	return nil
}

func (m *CmdsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "esc" && m.zone.Get("esc").InBounds(tea.MouseMsg{}) {
			return m, tea.Quit
		}
	case tea.MouseMsg:
		if m.zone.Get("esc").InBounds(msg) {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *CmdsModel) View() string {
	return m.zone.Mark("esc", "Click here or press Esc to go back")
}
