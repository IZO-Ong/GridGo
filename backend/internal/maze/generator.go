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
	Weight int // Priority value for Kruskal's
}

// initializeWallWeights sets every wall weight in the grid to a specific value.
// This is used to ensure solid colors for non-image-based generation modes.
func (m *Maze) initializeWallWeights(val int) {
	for r := 0; r < m.Rows; r++ {
		for c := 0; c < m.Cols; c++ {
			for i := 0; i < 4; i++ {
				m.Grid[r][c].WallWeights[i] = val
			}
		}
	}
}

// GenerateKruskal triggers a standard randomized Kruskal's generation.
func (m *Maze) GenerateKruskal() {
	m.initializeWallWeights(255)
	m.generateWeightedKruskal(nil)
}

// GenerateImageMaze triggers a guided Kruskal's generation using weights
// derived from an image.
func (m *Maze) GenerateImageMaze(weights map[string]int) {
	m.Weights = weights
	m.generateWeightedKruskal(weights)
}

// generateWeightedKruskal implements the core spanning tree logic.
func (m *Maze) generateWeightedKruskal(edgeWeights map[string]int) {
	dsu := NewDSU(m.Rows * m.Cols)
	var walls []Wall
	isImageMode := edgeWeights != nil

	for r := 0; r < m.Rows; r++ {
		for c := 0; c < m.Cols; c++ {
			if r < m.Rows-1 {
				w := Wall{R1: r, C1: c, R2: r + 1, C2: c}
				if val, ok := edgeWeights[fmt.Sprintf("%d-%d-top", r+1, c)]; ok {
					w.Weight = val
				} else {
					w.Weight = rand.IntN(100) 
				}
				walls = append(walls, w)
			}
			if c < m.Cols-1 {
				w := Wall{R1: r, C1: c, R2: r, C2: c + 1}
				if val, ok := edgeWeights[fmt.Sprintf("%d-%d-left", r, c+1)]; ok {
					w.Weight = val
				} else {
					w.Weight = rand.IntN(100)
				}
				walls = append(walls, w)
			}
		}
	}

	sort.Slice(walls, func(i, j int) bool {
		return walls[i].Weight < walls[j].Weight
	})

	for _, w := range walls {
		if isImageMode {
			if w.R1 == w.R2 { 
				m.Grid[w.R1][w.C1].WallWeights[1] = w.Weight 
				m.Grid[w.R2][w.C2].WallWeights[3] = w.Weight
			} else {
				m.Grid[w.R1][w.C1].WallWeights[2] = w.Weight
				m.Grid[w.R2][w.C2].WallWeights[0] = w.Weight
			}
		}

		id1 := w.R1*m.Cols + w.C1
		id2 := w.R2*m.Cols + w.C2

		if dsu.Find(id1) != dsu.Find(id2) {
			m.RemoveWalls(w.R1, w.C1, w.R2, w.C2)
			dsu.Union(id1, id2)
		}
	}
}

// GenerateRecursive sets up the grid with 255 weights before starting the DFS.
func (m *Maze) GenerateRecursive(r, c int) {
	// If this is the starting call, initialize the weights. 
	// (Note: in a recursive context, you might prefer calling initialize 
	// in a public wrapper, but this works if r/c are 0).
	if r == 0 && c == 0 {
		m.initializeWallWeights(255)
	}
	m.recursiveDFS(r, c)
}

func (m *Maze) recursiveDFS(r, c int) {
	m.Grid[r][c].Visited = true
	dirs := [][]int{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}
	rand.Shuffle(len(dirs), func(i, j int) {
		dirs[i], dirs[j] = dirs[j], dirs[i]
	})

	for _, d := range dirs {
		nextR, nextC := r+d[0], c+d[1]
		if nextR >= 0 && nextR < m.Rows && nextC >= 0 && nextC < m.Cols {
			if !m.Grid[nextR][nextC].Visited {
				m.RemoveWalls(r, c, nextR, nextC)
				m.recursiveDFS(nextR, nextC)
			}
		}
	}
}