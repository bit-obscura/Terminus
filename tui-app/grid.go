package main

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
}
