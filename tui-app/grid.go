package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type GridItem interface {
	View() string
}

type Grid struct {
	Rows    int
	Cols    int
	Items   [][]GridItem
	CursorX int
	CursorY int
}

func NewGrid(rows, cols int, items [][]GridItem) *Grid {
	return &Grid{
		Rows:    rows,
		Cols:    cols,
		Items:   items,
		CursorX: 0,
		CursorY: 0,
	}
}

func (g *Grid) MoveCursor(dx, dy int) {
	g.CursorX += dx
	g.CursorY += dy
	if g.CursorX < 0 {
		g.CursorX = 0
	} else if g.CursorX >= g.Cols {
		g.CursorX = g.Cols - 1
	}
	if g.CursorY < 0 {
		g.CursorY = 0
	} else if g.CursorY >= g.Rows {
		g.CursorY = g.Rows - 1
	}
}

func (g *Grid) View() string {
	var sb strings.Builder
	for y := 0; y < g.Rows; y++ {
		for x := 0; x < g.Cols; x++ {
			var cell string
			if y < len(g.Items) && x < len(g.Items[y]) && g.Items[y][x] != nil {
				cell = g.Items[y][x].View()
			} else {
				cell = ""
			}
			if x == g.CursorX && y == g.CursorY {
				cell = lipgloss.NewStyle().Bold(true).Render(cell)
			}
			if x > 0 {
				sb.WriteString(" | ")
			}
			sb.WriteString(cell)
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

type TextItem struct {
	Content string
	Style   lipgloss.Style
}

func (t TextItem) View() string {
	return t.Style.Render(t.Content)
}

type ListItem struct {
	Items []string
	Style lipgloss.Style
}

func (l ListItem) View() string {
	return l.Style.Render("- " + strings.Join(l.Items, "\n- "))
}

type TableItem struct {
	Data   [][]string
	Header lipgloss.Style
}

func (t TableItem) View() string {
	var sb strings.Builder
	for i, row := range t.Data {
		line := strings.Join(row, " \t ")
		if i == 0 {
			sb.WriteString(t.Header.Render(line) + "\n")
		} else {
			sb.WriteString(line + "\n")
		}
	}
	return sb.String()
}
