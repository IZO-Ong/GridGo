// Package maze handles the structural logic of the grid.
// This file specifically implements the generation algorithms.
package maze

import (
	"fmt"
	"math/rand/v2"
	"sort"
)

// Wall represents a potential boundary between two adjacent cells.
// Used as an edge in the graph representation of the maze.
type Wall struct {
	R1, C1 int // Coordinates of the first cell
	R2, C2 int // Coordinates of the second cell
	Weight int // Priority value; lower weights are more likely to be removed
}

// initializeWallWeights sets every wall weight in the grid to a specific value.
func (m *Maze) initializeWallWeights(val int) {
	for r := range m.Rows {
		for c := range m.Cols {
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

// GenerateImageMaze triggers a guided Kruskal's generation.
// Instead of random weights, it uses luminosity from an image 
// to decide which walls to break first.
func (m *Maze) GenerateImageMaze(weights map[string]int) {
	m.Weights = weights

	// map external border weights into internal grid structure
	for r := range m.Rows {
		for c := range m.Cols {
			if r == 0 {
				if w, ok := weights[fmt.Sprintf("%d-%d-top", r, c)]; ok {
					m.Grid[r][c].WallWeights[0] = w
				}
			}
			if c == 0 {
				if w, ok := weights[fmt.Sprintf("%d-%d-left", r, c)]; ok {
					m.Grid[r][c].WallWeights[3] = w
				}
			}
			if r == m.Rows-1 {
				if w, ok := weights[fmt.Sprintf("%d-%d-bottom", r, c)]; ok {
					m.Grid[r][c].WallWeights[2] = w
				}
			}
			if c == m.Cols-1 {
				if w, ok := weights[fmt.Sprintf("%d-%d-right", r, c)]; ok {
					m.Grid[r][c].WallWeights[1] = w
				}
			}
		}
	}

	m.generateWeightedKruskal(weights)
}

// generateWeightedKruskal implements the core spanning tree logic using a DSU
func (m *Maze) generateWeightedKruskal(edgeWeights map[string]int) {
	dsu := NewDSU(m.Rows * m.Cols)
	var walls []Wall
	isImageMode := edgeWeights != nil

	// Build a list of all internal walls
	for r := range m.Rows {
		for c := range m.Cols {
			if r < m.Rows-1 {
				w := Wall{R1: r, C1: c, R2: r + 1, C2: c}
				if val, ok := edgeWeights[fmt.Sprintf("%d-%d-top", r+1, c)]; ok {
					w.Weight = val
				} else {
					w.Weight = rand.IntN(100) // Random priority if no image weight
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

	// Sort walls by weight, with lowest weight walls are processed and removed first
	sort.Slice(walls, func(i, j int) bool {
		return walls[i].Weight < walls[j].Weight
	})

	for _, w := range walls {
		// Sync visual weights for rendering
		if isImageMode {
			if w.R1 == w.R2 { 
				m.Grid[w.R1][w.C1].WallWeights[1] = w.Weight 
				m.Grid[w.R2][w.C2].WallWeights[3] = w.Weight
			} else {
				m.Grid[w.R1][w.C1].WallWeights[2] = w.Weight
				m.Grid[w.R2][w.C2].WallWeights[0] = w.Weight
			}
		}

		// ID for DSU tracking
		id1 := w.R1*m.Cols + w.C1
		id2 := w.R2*m.Cols + w.C2

		// Only remove a wall if the two cells are not already connected.
		if dsu.Find(id1) != dsu.Find(id2) {
			m.RemoveWalls(w.R1, w.C1, w.R2, w.C2)
			dsu.Union(id1, id2)
		}
	}
}

// GenerateRecursive starts a Depth-First Search (DFS) generation.
// This results in mazes with long, winding paths and fewer junctions 
// compared to Kruskal's.
func (m *Maze) GenerateRecursive(r, c int) {
	if r == 0 && c == 0 {
		m.initializeWallWeights(255)
	}
	m.recursiveDFS(r, c)
}

// recursiveDFS is the internal engine for backtracking generation.
func (m *Maze) recursiveDFS(r, c int) {
	m.Grid[r][c].Visited = true
	dirs := [][]int{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}
	
	// Shuffle directions to ensure path does not favor one side
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