package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {

	model := NewRootModel()

	// WithAltScreen() is used to switch to the alternate screen buffer, which is a common practice in terminal applications.
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error starting the TUI:", err)
		os.Exit(1)
	}
}
