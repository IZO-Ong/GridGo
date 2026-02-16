package maze

import (
	"container/heap"
	"math"
)

type Point [2]int

// Priority Queue for A*
type Item struct {
	point    Point
	priority int
	index    int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].priority < pq[j].priority }
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index, pq[j].index = i, j
}
func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

// SolveAStar uses Manhattan Heuristic
func (m *Maze) SolveAStar() ([][2]int, [][2]int) {
	start, end := Point{m.Start[0], m.Start[1]}, Point{m.End[0], m.End[1]}
	visited, cameFrom, gScore := [][2]int{}, make(map[Point]Point), make(map[Point]int)
	gScore[start] = 0

	pq := &PriorityQueue{}
	heap.Init(pq)
	heap.Push(pq, &Item{point: start, priority: 0})

	for pq.Len() > 0 {
		curr := heap.Pop(pq).(*Item).point
		visited = append(visited, [2]int{curr[0], curr[1]})

		if curr == end { return visited, m.reconstructPath(cameFrom, curr) }

		for _, next := range m.GetNeighbors(curr) {
			tentativeG := gScore[curr] + 1
			if val, ok := gScore[next]; !ok || tentativeG < val {
				cameFrom[next] = curr
				gScore[next] = tentativeG
				fScore := tentativeG + m.manhattan(next, end)
				heap.Push(pq, &Item{point: next, priority: fScore})
			}
		}
	}
	return visited, nil
}

// SolveBFS for shortest path in unweighted grid
func (m *Maze) SolveBFS() ([][2]int, [][2]int) {
	start, end := Point{m.Start[0], m.Start[1]}, Point{m.End[0], m.End[1]}
	visited, queue, cameFrom := [][2]int{}, []Point{start}, make(map[Point]Point)
	seen := map[Point]bool{start: true}

	for len(queue) > 0 {
		curr := queue[0]; queue = queue[1:]
		visited = append(visited, [2]int{curr[0], curr[1]})

		if curr == end { return visited, m.reconstructPath(cameFrom, curr) }

		for _, next := range m.GetNeighbors(curr) {
			if !seen[next] {
				seen[next], cameFrom[next] = true, curr
				queue = append(queue, next)
			}
		}
	}
	return visited, nil
}

func (m *Maze) SolveGreedy() ([][2]int, [][2]int) {
	start, end := Point{m.Start[0], m.Start[1]}, Point{m.End[0], m.End[1]}
	visited, cameFrom := [][2]int{}, make(map[Point]Point)
	seen := map[Point]bool{start: true}

	pq := &PriorityQueue{}
	heap.Init(pq)
	// Initial priority is just the distance to the end
	heap.Push(pq, &Item{point: start, priority: m.manhattan(start, end)})

	for pq.Len() > 0 {
		curr := heap.Pop(pq).(*Item).point
		visited = append(visited, [2]int{curr[0], curr[1]})

		if curr == end {
			return visited, m.reconstructPath(cameFrom, curr)
		}

		for _, next := range m.GetNeighbors(curr) {
			if !seen[next] {
				seen[next] = true
				cameFrom[next] = curr
				priority := m.manhattan(next, end)
				heap.Push(pq, &Item{point: next, priority: priority})
			}
		}
	}
	return visited, [][2]int{}
}

func (m *Maze) manhattan(p1, p2 Point) int {
	return int(math.Abs(float64(p1[0]-p2[0])) + math.Abs(float64(p1[1]-p2[1])))
}

func (m *Maze) reconstructPath(cameFrom map[Point]Point, current Point) [][2]int {
	path := [][2]int{}
	for {
		path = append([][2]int{{current[0], current[1]}}, path...)
		if p, ok := cameFrom[current]; ok { current = p } else { break }
	}
	return path
}