package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Screen int

const (
	HomeScreen     Screen = iota // HomeScreen is the default screen
	SelectScreen                 // SelectScreen: the user can select an option between settings, apps, commands and explore
	ActivityScreen               // ActivityScreen: the user can see the activity of the selected option
)

type RootModel struct {
	CurrentScreen Screen
	Home          *HomeModel
	Select        *SelectModel
	// Activity      *ActivityModel
}

func NewRootModel() *RootModel {
	return &RootModel{
		CurrentScreen: HomeScreen,
		Home:          NewHomeModel(),
		Select:        NewSelectModel(),
		// Activity:      NewActivityModel(),
	}
}

// RootModel has to be passed by reference instead of by value.
// This is because each screen model contains sync/atomic.Bool from the bubblezone package.
// Since it's not allowed to copy a struct with a sync/atomic.Bool, we have to pass RootModel by reference.

func (r *RootModel) Init() tea.Cmd {
	return nil
}

func (r *RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	// Handle key events
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if r.CurrentScreen == HomeScreen {
				r.CurrentScreen = SelectScreen
			}
		case "esc":
			if r.CurrentScreen == SelectScreen {
				r.CurrentScreen = HomeScreen
			} else if r.CurrentScreen == ActivityScreen {
				r.CurrentScreen = SelectScreen
			}
		}
		// Handle mouse events
	case tea.MouseMsg:
		switch r.CurrentScreen {
		case HomeScreen:
			var m tea.Model
			m, cmd = r.Home.Update(msg)
			r.Home = m.(*HomeModel)
		case SelectScreen:
			var m tea.Model
			m, cmd = r.Select.Update(msg)
			r.Select = m.(*SelectModel)
			// case ActivityScreen:
			// 	var m tea.Model
			// 	m, cmd = r.Activity.Update(msg)
			// 	r.Activity = m.(*ActivityModel)
		}
	}

	return r, cmd
}

func (r *RootModel) View() string {
	switch r.CurrentScreen {
	case HomeScreen:
		return r.Home.View()
	case SelectScreen:
		return r.Select.View()
	// case ActivityScreen:
	// 	return r.Activity.View()
	default:
		return ""
	}
}
