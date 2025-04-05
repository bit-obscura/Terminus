package main

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ActionModel struct {
	CurrentScreen Screen
	blink         bool
	width         int
	height        int

	grid       Grid
	parentGrid Grid
	xPos       int
	yPos       int
}

func NewActionModel(grid Grid) *ActionModel {
	return &ActionModel{
		CurrentScreen: ActionScreen,
		grid:          grid,
	}
}

type Grid [][]interface{}

func (m ActionModel) Init() tea.Cmd {
	return tea.Tick(time.Second/2, func(time.Time) tea.Msg {
		return blinkMsg{}
	})
}

func (m *ActionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			m.yPos--
			if m.yPos < 0 {
				m.yPos = len(m.grid) - 1 // Wrap around to the bottom
			}
			m.xPos = 0 // Reset xPos when moving to a new row
		case "down":
			m.yPos++
			if m.yPos >= len(m.grid) {
				m.yPos = 0 // Wrap around to the top
			}
			m.xPos = 0 // Reset xPos when moving to a new row
		case "left":
			m.xPos--
			if m.xPos < 0 {
				m.xPos = len(m.grid[m.yPos]) - 1 // Wrap around to the last column in the row
			}
		case "right":
			m.xPos++
			if m.xPos >= len(m.grid[m.yPos]) {
				m.xPos = 0 // Wrap around to the first column in the row
			}
		case "enter":
			// Handle nested list navigation
			if nestedList, ok := m.grid[m.yPos][m.xPos].([]interface{}); ok {
				m.parentGrid = m.grid                // Save the current grid
				m.grid = [][]interface{}{nestedList} // Replace grid with the nested list
				m.xPos, m.yPos = 0, 0                // Reset position
			}
		case "backspace":
			// Handle going back to the parent grid
			if m.parentGrid != nil {
				m.grid = m.parentGrid
				m.parentGrid = nil
				m.xPos, m.yPos = 0, 0
			}
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m ActionModel) View() string {
	var output string

	for y, row := range m.grid {
		for x, cell := range row {
			cellStr := fmt.Sprintf("%v", cell)
			if x == m.xPos && y == m.yPos {
				output += lipgloss.NewStyle().
					Foreground(lipgloss.Color("black")).
					Background(lipgloss.Color("white")).
					Render(cellStr) + " "
			} else {
				output += cellStr + " "
			}
		}
		output += "\n"
	}

	return output
}
