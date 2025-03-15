package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	zone "github.com/lrstanley/bubblezone"
)

var activities = []string{
	"Settings",
	"Apps",
	"Commands",
	"Explore",
}

type SelectModel struct {
	Cursor int
	zone   zone.Manager
	id     string
}

func NewSelectModel() *SelectModel {
	return &SelectModel{
		id: zone.NewPrefix(),
	}
}

func (m *SelectModel) Init() tea.Cmd {
	return nil
}

func (m *SelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.Cursor > 0 {
				m.Cursor--
			}
		case "down":
			if m.Cursor < len(activities)-1 {
				m.Cursor++
			}
		case "enter":
			return m, tea.Quit
		}
	case tea.MouseMsg:
		for i := range activities {
			if m.zone.Get(fmt.Sprintf("activity-%d", i)).InBounds(msg) {
				m.Cursor = i
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m *SelectModel) View() string {
	view := "Select an activity:\n\n"
	for i, a := range activities {
		cursor := " "
		if i == m.Cursor {
			cursor = ">"
		}
		view += m.zone.Mark(fmt.Sprintf("activity-%d", i), fmt.Sprintf("%s %s\n", cursor, a))
	}
	return view
}
