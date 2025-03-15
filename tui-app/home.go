package main

import (
	tea "github.com/charmbracelet/bubbletea"

	zone "github.com/lrstanley/bubblezone"
)

type HomeModel struct {
	zone zone.Manager
	id   string
}

func NewHomeModel() *HomeModel {
	return &HomeModel{
		id: zone.NewPrefix(),
	}
}

func (m *HomeModel) Init() tea.Cmd {
	return nil
}

func (m *HomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.zone.Get("continue").InBounds(tea.MouseMsg{}) && msg.String() == "enter" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *HomeModel) View() string {
	return m.zone.Mark("continue", "Click here or press enter to continue")
}
