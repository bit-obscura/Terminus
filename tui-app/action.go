package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ActionModel struct {
	CurrentScreen Screen
	blink         bool
	width         int
	height        int

	xPos    int
	yPos    int
	tCursor []int
	int     // x, y
	active  bool
	id      int
	item    string
}

func NewActionModel() *ActionModel {
	return &ActionModel{
		CurrentScreen: HomeScreen,
	}
}

func (m ActionModel) Init() tea.Cmd {
	return tea.Tick(time.Second/2, func(time.Time) tea.Msg {
		return blinkMsg{}
	})
}

func (m ActionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m, nil
		case "up":
			m.yPos--
			if m.yPos < 0 {
				m.yPos = m.yPos.length
			}
			m.tCursor = append(m.tCursor, m.xPos, m.yPos)
			return m, nil
		case "down":
			m.yPos++
			if m.yPos > m.yPos.length {
				m.yPos = 0
			}
			m.tCursor = append(m.tCursor, m.xPos, m.yPos)
			return m, nil
		case "left":
			if m.xPos < 0 {
				m.xPos = m.xPos.length
			} else if m.xPos.length == 0 {
				m.yPos--
				if m.yPos < 0 {
					m.yPos = m.yPos.length
				}
				m.xPos = m.xPos.length
			} else {
				m.xPos--
			}
			m.tCursor = append(m.tCursor, m.xPos, m.yPos)
			return m, nil
		case "right":
			if m.xPos > m.xPos.length {
				m.xPos = 0
			} else if m.xPos.length == 0 {
				m.yPos++
				if m.yPos > m.yPos.length {
					m.yPos = 0
				}
				m.xPos = 0
			} else {
				m.xPos++
			}
			m.tCursor = append(m.tCursor, m.xPos, m.yPos)
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

func (m ActionModel) View() string {
	// Define the styles for the cursor and the text
	cursorStyle := lipgloss.NewStyle().Background(lipgloss.Color("#FF00FF")).Foreground(lipgloss.Color("#000000"))
	textStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#000000"))

	// Define the cursor position
	cursorPosition := m.tCursor[len(m.tCursor)-2:]

	// Define the text
	text := "Hello, World!"

	// Define the cursor
	cursor := "█"

	// Define the final text
	finalText := ""

	// Loop through the text
	for i, char := range text {
		// Check if the cursor is at the current position
		if cursorPosition[0] == i {
			// Add the cursor to the final text
			finalText += cursor
		}

		// Add the character to the final text
		finalText += string(char)
	}

	// Return the final text with the cursor
	return cursorStyle.Render(finalText)
}
