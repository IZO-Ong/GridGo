package main

import (
	"fmt"
	"math/rand"
)

type Cell struct {
	Row, Col int
	Visited  bool
	Walls    [4]bool
}

type Maze struct {
	Rows, Cols int
	Grid       [][]Cell
}

type DSU struct {
	parent []int
	rank   []int
}

type Wall struct {
	R1, C1 int
	R2, C2 int
}

func NewDSU(n int) *DSU {
	p := make([]int, n)
	r := make([]int, n)
	for i := range p {
		p[i] = i
		r[i] = 0
	}

	return &DSU{parent: p, rank: r}
}

func (d *DSU) Find(i int) int {
	if d.parent[i] == i {
		return i
	}
	d.parent[i] = d.Find(d.parent[i])
	return d.parent[i]
}

func (d *DSU) Union(i, j int) {
	rootI := d.Find(i)
	rootJ := d.Find(j)

	if rootI != rootJ {
		if d.rank[rootI] < d.rank[rootJ] {
			d.parent[rootI] = rootJ
		} else if d.rank[rootI] > d.rank[rootJ] {
			d.parent[rootJ] = rootI
		} else {
			d.parent[rootI] = rootJ
			d.rank[rootJ]++
		}
	}
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

func newMaze(rows, cols int) *Maze {
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

func main() {
	var rows, cols, choice int

	for {
		fmt.Print("Enter number of rows (minimum 2): ")
		_, err := fmt.Scan(&rows)

		if err == nil && rows >= 2 {
			break
		}

		fmt.Println("Invalid input, please try again.")
	}

	for {
		fmt.Print("Enter number of columns (minimum 2): ")
		_, err := fmt.Scan(&cols)

		if err == nil && cols >= 2 {
			break
		}

		fmt.Println("Invalid input, please try again.")
	}

	for {
		fmt.Println("\nChoose Generation Algorithm:")
		fmt.Println("1. Randomized Kruskal's (Short passages, many dead ends)")
		fmt.Println("2. Recursive Backtracker (Long, winding corridors)")
		fmt.Print("Selection: ")
		_, err := fmt.Scan(&choice)

		if err == nil && (choice == 1 || choice == 2) {
			break
		}
		fmt.Println("Invalid choice, please enter 1 or 2.")
	}

	maze := newMaze(rows, cols)

	if choice == 1 {
		maze.GenerateKruskal()
	} else {
		maze.GenerateRecursive(0, 0)
	}

	maze.Grid[0][0].Walls[0] = false
	maze.Grid[rows-1][cols-1].Walls[2] = false

	maze.Print()
}
