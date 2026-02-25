// Package maze provides the core data structures and logic for maze generation,
// spatial analysis, and pathfinding.
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
	ID         string         `json:"id"`
	Rows       int            `json:"rows"`
	Cols       int            `json:"cols"`
	Start      [2]int         `json:"start"`      // [row, col]
	End        [2]int         `json:"end"`        // [row, col]
	Grid       [][]Cell       `json:"grid"`       // The physical layout of cells and walls
	Weights    map[string]int `json:"weights"`    // Serialized weights for persistence/rendering
	Complexity float64        `json:"complexity"` // Calculated difficulty score
}

// Cell represents a single coordinate in the maze.
// It tracks its own boundaries and visual metadata for
// the rendering pipeline.
type Cell struct {
	Row         int     `json:"row"`
	Col         int     `json:"col"`
	Visited     bool    `json:"visited"`      // Used primarily during generation algorithms
	Walls       [4]bool `json:"walls"`        // Clockwise: 0:Top, 1:Right, 2:Bottom, 3:Left
	WallWeights [4]int  `json:"wall_weights"` // Visual/Difficulty weight of each wall
}

// MazeStats provides quantitative insights into the maze's topological properties.
type MazeStats struct {
	DeadEnds     int     `json:"dead_ends"`     // Cells with only 1 exit
	Junctions    int     `json:"junctions"`    // Cells with 3 or 4 exits (branching points)
	StraightWays int     `json:"straight_ways"` // Cells with 2 exits
	Complexity   float64 `json:"complexity"`    // Weighted score of branching vs size
}

// SetManualStartEnd allows specific placement of entrance/exit.
// Start and End points are automatically "clipped" (walls removed) to allow entry.
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
// It ensures a minimum Manhattan distance to prevent the start and end 
// from being too close to each other.
func (m *Maze) SetRandomStartEnd() {
	// Manhattan distance threshold: Ensure start/end are across at least 50% of the grid.
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

// getRandomBorderPoint returns coordinates on one of the four outer edges.
func (m *Maze) getRandomBorderPoint() (int, int) {
	side := rand.IntN(4)
	switch side {
	case 0: // Top edge
		return 0, rand.IntN(m.Cols)
	case 1: // Right edge
		return rand.IntN(m.Rows), m.Cols - 1
	case 2: // Bottom edge
		return m.Rows - 1, rand.IntN(m.Cols)
	default: // Left edge
		return rand.IntN(m.Rows), 0
	}
}

// clipBorderWall removes the outer-facing wall of a border cell to create an entrance/exit.
func (m *Maze) clipBorderWall(r, c int) {
	if r == 0 { m.Grid[r][c].Walls[0] = false }
	if r == m.Rows-1 { m.Grid[r][c].Walls[2] = false }
	if c == 0 { m.Grid[r][c].Walls[3] = false }
	if c == m.Cols-1 { m.Grid[r][c].Walls[1] = false }
}

// NewMaze initializes a grid where every cell is completely enclosed (all walls = true).
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

// RemoveWalls breaks the boundaries between two adjacent cells.
// It handles both horizontal and vertical adjacency.
func (m *Maze) RemoveWalls(r1, c1, r2, c2 int) {
	if r1 == r2 {
		// Horizontal neighbors
		if c1 < c2 {
			m.Grid[r1][c1].Walls[1] = false // Current Right
			m.Grid[r2][c2].Walls[3] = false // Target Left
		} else {
			m.Grid[r1][c1].Walls[3] = false // Current Left
			m.Grid[r2][c2].Walls[1] = false // Target Right
		}
	} else {
		// Vertical neighbors
		if r1 < r2 {
			m.Grid[r1][c1].Walls[2] = false // Current Bottom
			m.Grid[r2][c2].Walls[0] = false // Target Top
		} else {
			m.Grid[r1][c1].Walls[0] = false // Current Top
			m.Grid[r2][c2].Walls[2] = false // Target Bottom
		}
	}
}

// GetNeighbors returns a slice of adjacent points that can be reached.
// A point is a neighbor only if there is no wall between it and the current point.
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
			// Check the specific wall index for the direction of movement
			if !m.Grid[r][c].Walls[wallIdx] {
				neighbors = append(neighbors, Point{nr, nc})
			}
		}
	}

	return neighbors
}

// CalculateStats evaluates the maze's difficulty by analyzing junctions and dead ends.
// Complexity Heuristic: (Branching Factor * log2(TotalCells)).
func (m *Maze) CalculateStats() MazeStats {
	stats := MazeStats{}
	totalCells := float64(m.Rows * m.Cols)

	for r := range m.Rows {
		for c := range m.Cols {
			openCount := 0
			for _, isWall := range m.Grid[r][c].Walls {
				if !isWall { openCount++ }
			}

			switch openCount {
			case 1: stats.DeadEnds++
			case 2: stats.StraightWays++
			case 3, 4: stats.Junctions++
			}
		}
	}
	
	if totalCells > 0 {
		// Branching factor represents how many choices a user has on average.
		branchingFactor := (float64(stats.Junctions)*2.0 + float64(stats.DeadEnds)) / totalCells
		scaleBonus := math.Log2(totalCells)
		
		stats.Complexity = branchingFactor * scaleBonus
	}

	return stats
}

// SyncGridToWeights flattens the 2D grid wall data into a 1D map of string keys.
// This is used for JSON serialization and database storage.
func (m *Maze) SyncGridToWeights(original map[string]int) {
	m.Weights = make(map[string]int)
	for r := range m.Rows {
		for c := range m.Cols {
			keyTop := fmt.Sprintf("%d-%d-top", r, c)
			m.Weights[keyTop] = m.getWeightForWall(r, c, 0, keyTop, original)

			keyLeft := fmt.Sprintf("%d-%d-left", r, c)
			m.Weights[keyLeft] = m.getWeightForWall(r, c, 3, keyLeft, original)

			// Bottom and Right are only tracked for border edges to prevent redundancy
			if r == m.Rows-1 {
				keyBottom := fmt.Sprintf("%d-%d-bottom", r, c)
				m.Weights[keyBottom] = m.getWeightForWall(r, c, 2, keyBottom, original)
			}
			if c == m.Cols-1 {
				keyRight := fmt.Sprintf("%d-%d-right", r, c)
				m.Weights[keyRight] = m.getWeightForWall(r, c, 1, keyRight, original)
			}
		}
	}
}

// getWeightForWall determines the numeric weight of a wall (0 for paths, 120-255 for walls).
func (m *Maze) getWeightForWall(r, c, wallIdx int, key string, original map[string]int) int {
	if !m.Grid[r][c].Walls[wallIdx] {
		return 0 // Path (Weight 0)
	}
	if val, ok := original[key]; ok {
		return val
	}
	if original != nil {
		return 120 // Default wall weight for image-based mazes
	}
	return 255 // Absolute wall weight
}