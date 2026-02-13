package maze

import (
	"fmt"
	"math/rand/v2"
	"sort"
)

// Wall represents a potential boundary between two adjacent cells.
type Wall struct {
	R1, C1 int // Coordinates of the first cell
	R2, C2 int // Coordinates of the second cell
	Weight int // Priority value for Kruskal's; higher weights persist longer
}

// GenerateKruskal triggers a standard randomized Kruskal's generation.
func (m *Maze) GenerateKruskal() {
	m.generateWeightedKruskal(nil)
}

// GenerateImageMaze triggers a guided Kruskal's generation using weights
// derived from the Canny pipeline.
func (m *Maze) GenerateImageMaze(weights map[string]int) {
	m.generateWeightedKruskal(weights)
}

// generateWeightedKruskal implements the core spanning tree logic.
// It prioritizes removing walls with the lowest weights first, effectively
// forcing high-weight image edges to remain as part of the maze structure.
func (m *Maze) generateWeightedKruskal(edgeWeights map[string]int) {
	dsu := NewDSU(m.Rows * m.Cols)
	var walls []Wall

	// wall initalisation
	for r := range m.Rows {
		for c := range m.Cols {
			// horizontal separation
			if r < m.Rows-1 {
				w := Wall{R1: r, C1: c, R2: r + 1, C2: c}
				// map the bottom boundary of cell (r,c) to the "top" of (r+1,c)
				if val, ok := edgeWeights[fmt.Sprintf("%d-%d-top", r+1, c)]; ok {
					w.Weight = val
				} else {
					w.Weight = rand.IntN(100) // Random variance for filler paths
				}
				walls = append(walls, w)
			}
			if c < m.Cols-1 {
				w := Wall{R1: r, C1: c, R2: r, C2: c + 1}
				// map the right boundary of (r,c) to the "left" of (r,c+1)
				if val, ok := edgeWeights[fmt.Sprintf("%d-%d-left", r, c+1)]; ok {
					w.Weight = val
				} else {
					w.Weight = rand.IntN(100)
				}
				walls = append(walls, w)
			}
		}
	}

	// sort by weight
	sort.Slice(walls, func(i, j int) bool {
		return walls[i].Weight < walls[j].Weight
	})

	// kruskal algo
	for _, w := range walls {
		// persist weights into the grid for the renderer to use later
		if w.R1 == w.R2 { // vertical wall between columns
			m.Grid[w.R1][w.C1].WallWeights[1] = w.Weight // Right
			m.Grid[w.R2][w.C2].WallWeights[3] = w.Weight // Left
		} else { // horizontal wall between rows
			m.Grid[w.R1][w.C1].WallWeights[2] = w.Weight // Bottom
			m.Grid[w.R2][w.C2].WallWeights[0] = w.Weight // Top
		}

		id1 := w.R1*m.Cols + w.C1
		id2 := w.R2*m.Cols + w.C2

		if dsu.Find(id1) != dsu.Find(id2) {
			m.RemoveWalls(w.R1, w.C1, w.R2, w.C2)
			dsu.Union(id1, id2)
		}
	}
}

// GenerateRecursive implements a Depth-First Search (DFS) algorithm.
func (m *Maze) GenerateRecursive(r, c int) {
	m.Grid[r][c].Visited = true

	// randomly shuffle directions for unique corridors
	dirs := [][]int{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}
	rand.Shuffle(len(dirs), func(i, j int) {
		dirs[i], dirs[j] = dirs[j], dirs[i]
	})

	for _, d := range dirs {
		nextR, nextC := r+d[0], c+d[1]

		if nextR >= 0 && nextR < m.Rows && nextC >= 0 && nextC < m.Cols {
			if !m.Grid[nextR][nextC].Visited {
				m.RemoveWalls(r, c, nextR, nextC)
				m.GenerateRecursive(nextR, nextC)
			}
		}
	}
}
