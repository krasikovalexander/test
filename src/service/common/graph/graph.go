package graph

import (
	"container/list"
	"fmt"
)

//Graph is a data struct to store ordered graph with arbitrary node labels
type Graph struct {
	numNodes   int
	edges      [][]edge
	nodeLabels map[string]int
}

//Edge interface for edge value
type Edge interface {
	IsAccessibleFrom(edge interface{}) bool
}

//OptimalCriterion is an interface for optimization criterion
type OptimalCriterion interface {
	Apply(path *Path)
	GetResult() []*Path
}

type edge struct {
	from  int
	to    int
	value Edge
}

//Path represent slice of edges
type Path struct {
	edges []edge
}

//Edges get values of path edges
func (p *Path) Edges() []Edge {
	var edges []Edge
	for _, edge := range p.edges {
		edges = append(edges, edge.value)
	}
	return edges
}

//NewGraph creates graph with number of nodes required
func NewGraph(n int) *Graph {
	return &Graph{
		numNodes:   n,
		edges:      make([][]edge, n),
		nodeLabels: make(map[string]int),
	}
}

//AddEdge adds edges to graph. Creates nodes if they don't exist.
func (g *Graph) AddEdge(from string, to string, value Edge) {
	if from == to {
		return
	}

	u, exist := g.nodeLabels[from]
	if !exist {
		u = len(g.nodeLabels)
		g.nodeLabels[from] = u
	}

	v, exist := g.nodeLabels[to]
	if !exist {
		v = len(g.nodeLabels)
		g.nodeLabels[to] = v
	}

	g.edges[u] = append(g.edges[u], edge{from: u, to: v, value: value})
}

//GetPaths search paths between two nodes. Pass limit greater than zero to set maximim path length
func (g *Graph) GetPaths(from string, to string, limit int) []Path {
	var result []Path

	fromIdx, fromExists := g.nodeLabels[from]
	toIdx, toExists := g.nodeLabels[to]

	if from == to || !fromExists || !toExists {
		return result
	}

	paths := g.getPaths(fromIdx, toIdx, limit)

	for _, edges := range paths {
		result = append(result, Path{
			edges: edges,
		})
	}
	return result
}

//based on BFS algorithm
func (g *Graph) getPaths(from int, to int, limit int) [][]edge {

	type queueItem struct {
		path    []edge
		visited []bool
	}

	var path []edge
	visited := make([]bool, g.numNodes)
	visited[from] = true

	queue := list.New()
	for _, edge := range g.edges[from] {
		queue.PushBack(queueItem{
			path:    append(path, edge),
			visited: visited,
		})
	}

	var paths [][]edge

	for {
		next := queue.Front()
		if next == nil {
			break
		}
		queue.Remove(next)

		item := next.Value.(queueItem)

		pathLength := len(item.path)
		if pathLength == 0 || (limit > 0 && pathLength == limit) {
			continue
		}

		currentEdge := &item.path[len(item.path)-1]

		if currentEdge.to == to {
			paths = append(paths, item.path)
			continue
		}

		item.visited[currentEdge.to] = true

		for _, edge := range g.edges[currentEdge.to] {
			if !visited[edge.to] && edge.value.IsAccessibleFrom(currentEdge.value) {
				queue.PushBack(queueItem{
					path:    append(item.path, edge),
					visited: item.visited,
				})
			}
		}
	}
	return paths
}

//Print output simple debug info
func (g *Graph) Print() {
	fmt.Println(g.numNodes)
	fmt.Println(g.nodeLabels)
	fmt.Println(g.edges)
}

//SearchOptimalPaths search optimal paths between two nodes by given criteria
func (g *Graph) SearchOptimalPaths(from string, to string, limit int, criteria ...OptimalCriterion) {
	fromIdx, fromExists := g.nodeLabels[from]
	toIdx, toExists := g.nodeLabels[to]

	if from == to || !fromExists || !toExists {
		return
	}

	g.searchOptimalPaths(fromIdx, toIdx, limit, criteria...)
}

func (g *Graph) searchOptimalPaths(from int, to int, limit int, criteria ...OptimalCriterion) {
	type queueItem struct {
		path    []edge
		visited []bool
	}

	var path []edge
	visited := make([]bool, g.numNodes)
	visited[from] = true

	queue := list.New()
	for _, edge := range g.edges[from] {
		queue.PushBack(queueItem{
			path:    append(path, edge),
			visited: visited,
		})
	}

	for {
		next := queue.Front()
		if next == nil {
			break
		}
		queue.Remove(next)

		item := next.Value.(queueItem)

		pathLength := len(item.path)
		if pathLength == 0 || (limit > 0 && pathLength == limit) {
			continue
		}

		currentEdge := &item.path[len(item.path)-1]

		if currentEdge.to == to {
			for _, criterion := range criteria {
				criterion.Apply(&Path{edges: item.path})
			}
			continue
		}

		item.visited[currentEdge.to] = true

		for _, edge := range g.edges[currentEdge.to] {
			if !visited[edge.to] && edge.value.IsAccessibleFrom(currentEdge.value) {
				queue.PushBack(queueItem{
					path:    append(item.path, edge),
					visited: item.visited,
				})
			}
		}
	}
}
