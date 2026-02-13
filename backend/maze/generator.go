package maze

import (
	"fmt"
	"math/rand/v2"
	"sort"
)

type Wall struct {
	R1, C1 int
	R2, C2 int
	Weight int
}

func (m *Maze) GenerateKruskal() {
	m.generateWeightedKruskal(nil)
}

func (m *Maze) GenerateImageMaze(weights map[string]int) {
	m.generateWeightedKruskal(weights)
}

func (m *Maze) generateWeightedKruskal(edgeWeights map[string]int) {
	dsu := NewDSU(m.Rows * m.Cols)
	var walls []Wall

	for r := range m.Rows {
		for c := range m.Cols {
			if r < m.Rows-1 {
				w := Wall{R1: r, C1: c, R2: r + 1, C2: c}
				val, ok := edgeWeights[fmt.Sprintf("%d-%d-bottom", r, c)]

				if ok {
					w.Weight = val
				} else {
					w.Weight = rand.IntN(100) // random weight
				}
				walls = append(walls, w)
			}
			if c < m.Cols-1 {
				w := Wall{R1: r, C1: c, R2: r, C2: c + 1}
				val, ok := edgeWeights[fmt.Sprintf("%d-%d-right", r, c)]
				if ok {
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
		id1 := w.R1*m.Cols + w.C1
		id2 := w.R2*m.Cols + w.C2

		if dsu.Find(id1) != dsu.Find(id2) {
			m.RemoveWalls(w.R1, w.C1, w.R2, w.C2)
			dsu.Union(id1, id2)
		}
	}
}

func (m *Maze) GenerateRecursive(r, c int) {
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
				m.GenerateRecursive(nextR, nextC)
			}
		}
	}
}
