package maze

import "fmt"

type Maze struct {
	Rows, Cols int
	Grid       [][]Cell
}

type Cell struct {
	Row, Col int
	Visited  bool
	Walls    [4]bool
}

type Wall struct {
	R1, C1 int
	R2, C2 int
}

func NewMaze(rows, cols int) *Maze {
	grid := make([][]Cell, rows)

	for r := range rows {
		grid[r] = make([]Cell, cols)

		for c := range cols {
			grid[r][c] = Cell{
				Row:   r,
				Col:   c,
				Walls: [4]bool{true, true, true, true},
			}
		}

	}

	return &Maze{Rows: rows, Cols: cols, Grid: grid}
}

func (m *Maze) Print() {
	for r := range m.Rows {
		for c := range m.Cols {
			if m.Grid[r][c].Walls[0] {
				fmt.Print("+---")
			} else {
				fmt.Print("+   ")
			}
		}
		fmt.Println("+")

		for c := range m.Cols {
			if m.Grid[r][c].Walls[3] {
				fmt.Print("|   ")
			} else {
				fmt.Print("    ")
			}
		}
		fmt.Println("|")
	}

	for c := range m.Cols {
		if m.Grid[m.Rows-1][c].Walls[2] {
			fmt.Print("+---")
		} else {
			fmt.Print("+   ")
		}
	}
	fmt.Println("+")
}

func (m *Maze) RemoveWalls(r1, c1, r2, c2 int) {
	if r1 == r2 {
		if c1 < c2 {
			m.Grid[r1][c1].Walls[1] = false
			m.Grid[r2][c2].Walls[3] = false
		} else {
			m.Grid[r1][c1].Walls[3] = false
			m.Grid[r2][c2].Walls[1] = false
		}
	} else {
		if r1 < r2 {
			m.Grid[r1][c1].Walls[2] = false
			m.Grid[r2][c2].Walls[0] = false
		} else {
			m.Grid[r1][c1].Walls[0] = false
			m.Grid[r2][c2].Walls[2] = false
		}
	}
}
