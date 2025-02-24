package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	choises  []string         //item list
	cursor   int              //cursor position relative to the item list
	selected map[int]struct{} //selected items
}

func initModel() model {
	return model{
		//is a list of cities already visited
		choises: []string{"New York", "London", "Paris", "Tokyo", "Sydney", "Cape Town"},

		//Uso la mappa per tenere traccia degli elementi selezionati,
		//le chiavi sono gli indici degli elementi selezionati.
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	//No need I/O operations now
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	//Handle key presses
	case tea.KeyMsg:

		//Which key was pressed?
		switch msg.String() {
		case "ctrl+c", "q": //Quit the program
			return m, tea.Quit

		case "up", "k": //Move the cursor up
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j": //Move the cursor down
			if m.cursor < len(m.choises)-1 {
				m.cursor++
			}

		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	//Header of the list
	s := "Which cities have you visited?\n\n"

	//Iterate over the list of cities
	for i, city := range m.choises {

		//Cursor position
		cursor := " " //No cursor
		if m.cursor == i {
			cursor = ">" //Cursor
		}

		//Check if the city is selected
		checked := " " //Not selected
		if _, ok := m.selected[i]; ok {
			checked = "x" //Selected
		}

		//Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, city)
	}

	//Footer
	s += "\nPress 'q' to quit.\n"

	//Send the UI for rendering
	return s
}

func main() {
	p := tea.NewProgram(initModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("There's been an error: %v", err)
		os.Exit(1)
	}
}
