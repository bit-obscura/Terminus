package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	model := NewRootModel()
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error starting the TUI:", err)
	}
}
