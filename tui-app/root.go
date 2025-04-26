package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/gamut"
)

var (

	// General
	normal = lipgloss.Color("EEEEEE")
	hidden = lipgloss.Color("575755")
	teal   = lipgloss.Color("#006883")
	lblue  = lipgloss.Color("#31D8EE")
	green  = lipgloss.Color("#01DC94")
	yellow = lipgloss.Color("#DCFE54")

	// Styles
	base   = lipgloss.NewStyle().Foreground(normal)
	blends = gamut.Blends(lipgloss.Color(teal), lipgloss.Color(lblue), 50)

	titleStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color(lblue)).Bold(true)
	subtitleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(teal)).Italic(true)
	fakeBtnStyle  = lipgloss.NewStyle().
			Foreground(lipgloss.Color(green)).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(green)).
			Padding(0, 1).
			Align(lipgloss.Center)

	// Borders
	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(normal)).
			Margin(1, 2).
			Padding(2, 4)
)

type Screen int

const (
	HomeScreen Screen = iota
	ActionScreen
)

type RootModel struct {
	CurrentScreen Screen
	blink         bool
	width         int
	height        int
}

func NewRootModel() *RootModel {
	return &RootModel{
		CurrentScreen: HomeScreen,
	}
}

type blinkMsg struct{}

func createContainer(content string, width, height int) string {
	containerWidth := width - 6
	containerHeight := height - 4

	container := borderStyle.
		Width(containerWidth).
		Height(containerHeight).
		Align(lipgloss.Center, lipgloss.Center)

	return container.Render(content)
}

func (m *RootModel) Init() tea.Cmd {
	return tea.Tick(time.Second/2, func(time.Time) tea.Msg {
		return blinkMsg{}
	})
}

func (m *RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.CurrentScreen == HomeScreen {
				m.CurrentScreen = ActionScreen
			}
			return m, nil
		case "esc":
			if m.CurrentScreen == ActionScreen {
				m.CurrentScreen = HomeScreen
			}
			return m, nil
		case "ctrl+c":
			return m, tea.Quit
		}
	case blinkMsg:
		m.blink = !m.blink
		return m, tea.Tick(time.Second/2, func(time.Time) tea.Msg {
			return blinkMsg{}
		})
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

	return m, nil
}

func (m *RootModel) View() string {
	viewportWidth, viewportHeight := 120, 30 // Default size

	if m.width > 0 && m.height > 0 {
		viewportWidth, viewportHeight = m.width, m.height
	}

	switch m.CurrentScreen {
	case HomeScreen:
		title := titleStyle.Render("TUI")
		subtitle := subtitleStyle.Render("powered by CodeLab")

		var btn string
		if m.blink {
			btn = fakeBtnStyle.Render("Enter")
		} else {
			btn = fakeBtnStyle.Copy().Border(lipgloss.HiddenBorder()).Render(" ")
		}

		footer := lipgloss.NewStyle().Foreground(lipgloss.Color(hidden)).Render("Press 'Enter' to continue or 'Ctrl+c' to quit")

		content := title + "\n" + subtitle + "\n\n" + btn + "\n\n" + footer

		return createContainer(content, viewportWidth, viewportHeight)
	case ActionScreen:
		return NewActionModel().View()
	}

	return ""
}
