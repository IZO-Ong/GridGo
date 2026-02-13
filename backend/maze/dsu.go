package maze

// DSU (Disjoint Set Union) provides efficient tracking of
// connected components in the maze grid for kruskal algorithm
type DSU struct {
	parent []int
	rank   []int
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
	// compress path
	d.parent[i] = d.Find(d.parent[i])
	return d.parent[i]
}

func (d *DSU) Union(i, j int) {
	rootI := d.Find(i)
	rootJ := d.Find(j)

	if rootI != rootJ {
		// optimise using rank
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
