package graph

import (
	"container/list"
	"fmt"
)

type Graph struct {
	numNodes          int
	edges             [][]edge
	nodeLabels        map[string]int
	reverseNodeLabels []string
}

type Edge interface {
	IsAccessibleFrom(edge interface{}) bool
	Cost() float32
}

type edge struct {
	from  int
	to    int
	value Edge
}

type Path struct {
	edges []edge
}

func (p *Path) Edges() []Edge {
	var edges []Edge
	for _, edge := range p.edges {
		edges = append(edges, edge.value)
	}
	return edges
}

func NewGraph(n int) *Graph {
	return &Graph{
		numNodes:          n,
		edges:             make([][]edge, n),
		nodeLabels:        make(map[string]int),
		reverseNodeLabels: make([]string, n),
	}
}

func (g *Graph) AddEdge(from string, to string, value Edge) {
	if from == to {
		return
	}

	u, exist := g.nodeLabels[from]
	if !exist {
		u = len(g.nodeLabels)
		g.nodeLabels[from] = u
		g.reverseNodeLabels[u] = from
	}

	v, exist := g.nodeLabels[to]
	if !exist {
		v = len(g.nodeLabels)
		g.nodeLabels[to] = v
		g.reverseNodeLabels[v] = to
	}

	g.edges[u] = append(g.edges[u], edge{from: u, to: v, value: value})
}

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
		currentEdge := &item.path[len(item.path)-1]

		if currentEdge.to == to {
			paths = append(paths, item.path)
			continue
		}

		if limit > 0 && len(item.path) == limit {
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

func (g *Graph) Print() {
	fmt.Println(g.numNodes)
	fmt.Println(g.nodeLabels)
	fmt.Println(g.reverseNodeLabels)
	fmt.Println(g.edges)
}
