package maze

import (
	"fmt"
	"math"
	"math/rand/v2"
)

// Maze represents the overall grid structure.
// It serves as the primary state container for both generation
// and rendering logic.
type Maze struct {
    Rows    int            `json:"rows"`
    Cols    int            `json:"cols"`
    Grid    [][]Cell       `json:"grid"`
    Start   [2]int         `json:"start"`
    End     [2]int         `json:"end"`
    Weights map[string]int `json:"weights"` 
}

// Cell represents a single coordinate in the maze.
// It tracks its own boundaries and visual metadata for
// the rendering pipeline.
type Cell struct {
	Row         int     `json:"row"`
	Col         int     `json:"col"`
	Visited     bool    `json:"visited"`
	Walls       [4]bool `json:"walls"`        // 0:Top, 1:Right, 2:Bottom, 3:Left
	WallWeights [4]int  `json:"wall_weights"`
}

// SetManualStartEnd allows specific placement of entrance/exit
func (m *Maze) SetManualStartEnd(sr, sc, er, ec int) error {
	isBorder := func(r, c int) bool {
		return r == 0 || r == m.Rows-1 || c == 0 || c == m.Cols-1
	}

	if !isBorder(sr, sc) || !isBorder(er, ec) {
		return fmt.Errorf("start and end points must be on the maze border")
	}

	m.Start = [2]int{sr, sc}
	m.End = [2]int{er, ec}

	m.clipBorderWall(sr, sc)
	m.clipBorderWall(er, ec)
	return nil
}

// SetRandomStartEnd picks two unique points on the maze boundary.
func (m *Maze) SetRandomStartEnd() {
	// Manhattan distance threshold (at least 50% of max)
	minDist := float64(m.Rows+m.Cols) * 0.5

	for {
		sR, sC := m.getRandomBorderPoint()
		eR, eC := m.getRandomBorderPoint()

		dist := math.Abs(float64(sR-eR)) + math.Abs(float64(sC-eC))
		if (sR != eR || sC != eC) && dist >= minDist {
			m.Start = [2]int{sR, sC}
			m.End = [2]int{eR, eC}
			break
		}
	}

	m.clipBorderWall(m.Start[0], m.Start[1])
	m.clipBorderWall(m.End[0], m.End[1])
}

func (m *Maze) getRandomBorderPoint() (int, int) {
	side := rand.IntN(4)
	switch side {
	case 0:
		return 0, rand.IntN(m.Cols)
	case 1:
		return rand.IntN(m.Rows), m.Cols - 1
	case 2:
		return m.Rows - 1, rand.IntN(m.Cols)
	default:
		return rand.IntN(m.Rows), 0
	}
}

func (m *Maze) clipBorderWall(r, c int) {
	if r == 0 {
		m.Grid[r][c].Walls[0] = false
	}
	if r == m.Rows-1 {
		m.Grid[r][c].Walls[2] = false
	}
	if c == 0 {
		m.Grid[r][c].Walls[3] = false
	}
	if c == m.Cols-1 {
		m.Grid[r][c].Walls[1] = false
	}
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

		for c := range m.Cols {
			if m.Grid[r][c].Walls[3] {
				fmt.Print("|")
			} else {
				fmt.Print(" ")
			}

			// markers
			if r == m.Start[0] && c == m.Start[1] {
				fmt.Print(" * ")
			} else if r == m.End[0] && c == m.End[1] {
				fmt.Print(" - ")
			} else {
				fmt.Print("   ")
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

// GetNeighbors returns a slice of adjacent points that can be reached from 
// the current point (i.e., they are within bounds and not blocked by a wall).
func (m *Maze) GetNeighbors(p Point) []Point {
	neighbors := []Point{}
	r, c := p[0], p[1]

	dirs := [][]int{
		{-1, 0, 0}, // North
		{0, 1, 1},  // East
		{1, 0, 2},  // South
		{0, -1, 3}, // West
	}

	for _, d := range dirs {
		nr, nc := r+d[0], c+d[1]
		wallIdx := d[2]

		if nr >= 0 && nr < m.Rows && nc >= 0 && nc < m.Cols {
			if !m.Grid[r][c].Walls[wallIdx] {
				neighbors = append(neighbors, Point{nr, nc})
			}
		}
	}

	return neighbors
}