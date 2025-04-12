package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Items for the grid interface

// Text based item

type TextItem struct {
	Title   string
	Content string
	ID      string
	Style   lipgloss.Style
}

type TextInput struct {
	Title   string
	Content string
	ID      string
	Style   lipgloss.Style
}

type ListItem struct {
	Title    string
	Items    []string
	Cursor   int
	Selected map[int]struct{}
	ID       string
	Style    lipgloss.Style
}

type ButtonItem struct {
	Title   string
	Content string
	ID      string
	Style   lipgloss.Style
}

type TableItem struct {
	Title   string
	Headers []string
	Rows    [][]string
	ID      string
	Style   lipgloss.Style
}

func (t TextItem) View() string {
	return lipgloss.NewStyle().Render(t.Content)
}

func (t TextInput) View() string {
	return lipgloss.NewStyle().Render(t.Content)
}

func (l ListItem) View() string {
	var sb strings.Builder
	sb.WriteString(l.Title + "\n")
	for i, item := range l.Items {
		if _, ok := l.Selected[i]; ok {
			sb.WriteString("-> " + item + "\n")
		} else {
			sb.WriteString("   " + item + "\n")
		}
	}
	return sb.String()
}

func (b ButtonItem) View() string {
	return lipgloss.NewStyle().Render(b.Content)
}

func (t TableItem) View() string {
	var sb strings.Builder
	sb.WriteString(t.Title + "\n")
	for _, header := range t.Headers {
		sb.WriteString(header + "\t")
	}
	sb.WriteString("\n")
	for _, row := range t.Rows {
		for _, cell := range row {
			sb.WriteString(cell + "\t")
		}
		sb.WriteString("\n")
	}
	return sb.String()
}
