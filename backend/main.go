package main

import (
	"fmt"
	"math/rand"
	"strings"
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

func (m *Maze) getRandomNeighbours(r, c int) [][]int {
	neighbours := [][]int{}

	directions := [][]int{
		{-1, 0},
		{1, 0},
		{0, 1},
		{0, -1},
	}

	for _, d := range directions {
		nextR, nextC := r+d[0], c+d[1]

		if nextR >= 0 && nextR < m.Rows && nextC >= 0 && nextC < m.Cols {
			if !m.Grid[nextR][nextC].Visited {
				neighbours = append(neighbours, []int{nextR, nextC})
			}
		}
	}

	rand.Shuffle(len(neighbours), func(i, j int) {
		neighbours[i], neighbours[j] = neighbours[j], neighbours[i]
	})

	return neighbours
}

func (m *Maze) RemoveWalls(r1, c1, r2, c2 int) {
	if r1 == r2 {
		if c1 < c2 { // Move right
			m.Grid[r1][c1].Walls[1] = false
			m.Grid[r2][c2].Walls[3] = false
		} else { // Move Left
			m.Grid[r1][c1].Walls[3] = false
			m.Grid[r2][c2].Walls[1] = false
		}
	} else {
		if r1 < r2 {
			m.Grid[r1][c1].Walls[0] = false
			m.Grid[r2][c2].Walls[2] = false
		} else {
			m.Grid[r1][c1].Walls[2] = false
			m.Grid[r2][c2].Walls[0] = false
		}
	}
}

func (m *Maze) Generate(r, c int) {
	m.Grid[r][c].Visited = true

	neighbours := m.getRandomNeighbours(r, c)

	for _, coords := range neighbours {
		nextR, nextC := coords[0], coords[1]
		if !m.Grid[nextR][nextC].Visited {
			m.RemoveWalls(r, c, nextR, nextC)
			m.Generate(nextR, nextC)
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
	fmt.Println(strings.Repeat("+---", m.Cols) + "+")
}

func main() {
	maze := newMaze(10, 10)
	maze.Generate(0, 0)
	maze.Print()
}
