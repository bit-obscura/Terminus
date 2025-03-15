package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	zone "github.com/lrstanley/bubblezone"
)

func main() {

	defer zone.Close() // Close Bubblezone when the program exits
	model := NewRootModel()

	// WithAltScreen() is used to switch to the alternate screen buffer, which is a common practice in terminal applications.
	// WithMouseAllMotion() is used to enable mouse tracking in the terminal.
	p := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseAllMotion())

	if err := p.Start(); err != nil {
		fmt.Println("Error starting the TUI:", err)
	}
}
