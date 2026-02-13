package maze

import "math/rand/v2"

func (m *Maze) GenerateKruskal() {
	dsu := NewDSU(m.Rows * m.Cols)
	walls := []Wall{}

	for r := range m.Rows {
		for c := range m.Cols {
			if r < m.Rows-1 {
				walls = append(walls, Wall{r, c, r + 1, c})
			}
			if c < m.Cols-1 {
				walls = append(walls, Wall{r, c, r, c + 1})
			}
		}
	}

	rand.Shuffle(len(walls), func(i, j int) {
		walls[i], walls[j] = walls[j], walls[i]
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
