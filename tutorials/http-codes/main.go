package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const url = "https://www.google.com/maps"

type model struct {
	status int
	err    error
}

// command to make a request to the server and return the result as a msg.
func checkServer() tea.Msg {
	c := &http.Client{Timeout: 10 * time.Second}
	res, err := c.Get(url)

	if err != nil {
		return errMsg{err}
	}

	return statusMsg(res.StatusCode)
}

type (
	statusMsg int
	errMsg    struct{ err error }
)

func (e errMsg) Error() string { return e.err.Error() } // implement error interface

func (m model) Init() tea.Cmd {
	return checkServer
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case statusMsg:
		m.status = int(msg)
		return m, tea.Quit
	case errMsg:
		m.err = msg.err
		return m, tea.Quit
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("\nError: %v\n\n", m.err)
	}

	s := fmt.Sprintf("Checking %s...\n\n", url)

	if m.status > 0 {
		s += fmt.Sprintf("%d %s!", m.status, http.StatusText(m.status))
	}

	return "\n" + s + "\n\n"
}

func main() {
	if err := tea.NewProgram(model{}).Start(); err != nil {
		fmt.Printf("Error starting program: %v", err)
		os.Exit(1)
	}
}
