package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

var p *tea.Program

type SessionState int

const (
	MainView SessionState = iota
	ProjectView
	SettingsView
	ApplicationsView
	CmdGoalsView
	ExploreView
)

type MainModel struct {
	CurrentState SessionState
	ProjectModel *ProjectModel
	// SettingsModel tea.Model
	// ApplicationsModel tea.Model
	// CmdGoalsModel tea.Model
	// ExploreModel tea.Model
}

func NewMainModel() *MainModel {
	return &MainModel{
		CurrentState: MainView,
		ProjectModel: NewProjectModel(),
	}
}

func (m *MainModel) Init() tea.Cmd {
	return nil
}

func (m *MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.CurrentState == MainView {
			switch msg.String() {
			case "enter":
				m.CurrentState = ProjectView
				return m, nil
			case "ctrl+c", "q":
				return m, tea.Quit
			}
		} else if m.CurrentState == ProjectView {
			switch msg.String() {
			case "esc":
				m.CurrentState = MainView
				return m, nil
			}
			var cmd tea.Cmd
			updatedModel, cmd := m.ProjectModel.Update(msg)
			if updatedProjectModel, ok := updatedModel.(*ProjectModel); ok {
				m.ProjectModel = updatedProjectModel
			}
			return m, cmd
		}
	}

	return m, nil
}

func (m *MainModel) View() string {
	if m.CurrentState == MainView {
		return "Main View\n\nPress Enter to go to Project View\nPress Ctrl+C to exit\n"
	} else if m.CurrentState == ProjectView {
		return m.ProjectModel.View()
	}
	return ""
}

func main() {
	p = tea.NewProgram(NewMainModel(), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Printf("Error starting program: %v\n", err)
	}
}
