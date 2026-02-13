package maze

import "fmt"

// Maze represents the overall grid structure.
// It serves as the primary state container for both generation
// and rendering logic.
type Maze struct {
	Rows, Cols int
	Grid       [][]Cell
}

// Cell represents a single coordinate in the maze.
// It tracks its own boundaries and visual metadata for
// the rendering pipeline.
type Cell struct {
	Row, Col    int
	Visited     bool    // Used by DFS (Recursive Backtracker) algorithm
	Walls       [4]bool // 0:Top, 1:Right, 2:Bottom, 3:Left
	WallWeights [4]int  // Maps Canny magnitudes to individual wall intensity
}

// NewMaze initializes a grid where every cell is completely enclosed.
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

// Print outputs a rough ASCII representation of the maze to the terminal.
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

		// draw left/right walls
		for c := range m.Cols {
			if m.Grid[r][c].Walls[3] {
				fmt.Print("|   ")
			} else {
				fmt.Print("    ")
			}
		}
		fmt.Println("|")
	}

	// closing bottom edge for grid
	for c := range m.Cols {
		if m.Grid[m.Rows-1][c].Walls[2] {
			fmt.Print("+---")
		} else {
			fmt.Print("+   ")
		}
	}
	fmt.Println("+")
}

// RemoveWalls breaks the boundaries between two adjacent cells.
func (m *Maze) RemoveWalls(r1, c1, r2, c2 int) {
	if r1 == r2 {
		// Horizontal neighbors
		if c1 < c2 {
			m.Grid[r1][c1].Walls[1] = false // Right
			m.Grid[r2][c2].Walls[3] = false // Left
		} else {
			m.Grid[r1][c1].Walls[3] = false // Left
			m.Grid[r2][c2].Walls[1] = false // Right
		}
	} else {
		// Vertical neighbors
		if r1 < r2 {
			m.Grid[r1][c1].Walls[2] = false // Bottom
			m.Grid[r2][c2].Walls[0] = false // Top
		} else {
			m.Grid[r1][c1].Walls[0] = false // Top
			m.Grid[r2][c2].Walls[2] = false // Bottom
		}
	}
}
