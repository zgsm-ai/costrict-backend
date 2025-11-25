package main

import (
	"fmt"
	"math"
)

type Graph struct {
	vertices int
	adjList  [][]int
}

type Edge struct {
	from int
	to   int
	cost float64
}

func NewGraph(vertices int) *Graph {
	return &Graph{
		vertices: vertices,
		adjList:  make([][]int, vertices),
	}
}

func (g *Graph) AddEdge(u, v int) {
	g.adjList[u] = append(g.adjList[u], v)
	g.adjList[v] = append(g.adjList[v], u)
}

func (g *Graph) DFS(start int, visited []bool) {
	visited[start] = true
	fmt.Printf("%d ", start)
	
	for _, neighbor := range g.adjList[start] {
		if !visited[neighbor] {
			g.DFS(neighbor, visited)
		}
	}
}

func (g *Graph) BFS(start int) {
	visited := make([]bool, g.vertices)
	queue := []int{start}
	visited[start] = true
	
	for len(queue) > 0 {
		vertex := queue[0]
		queue = queue[1:]
		
		fmt.Printf("%d ", vertex)
		
		for _, neighbor := range g.adjList[vertex] {
			if !visited[neighbor] {
				visited[neighbor] = true
				queue = append(queue, neighbor)
			}
		}
	}
}

type WeightedGraph struct {
	vertices int
	edges    []Edge
}

func NewWeightedGraph(vertices int) *WeightedGraph {
	return &WeightedGraph{
		vertices: vertices,
		edges:    make([]Edge, 0),
	}
}

func (wg *WeightedGraph) AddWeightedEdge(u, v int, cost float64) {
	wg.edges = append(wg.edges, Edge{u, v, cost})
}

func (wg *WeightedGraph) find(parent []int, i int) int {
	if parent[i] == -1 {
		return i
	}
	return wg.find(parent, parent[i])
}

func (wg *WeightedGraph) union(parent []int, x, y int) {
	xSet := wg.find(parent, x)
	ySet := wg.find(parent, y)
	parent[xSet] = ySet
}

func (wg *WeightedGraph) KruskalMST() []Edge {
	result := make([]Edge, 0)
	
	// Sort all edges in non-decreasing order of their weight
	for i := 0; i < len(wg.edges)-1; i++ {
		for j := 0; j < len(wg.edges)-i-1; j++ {
			if wg.edges[j].cost > wg.edges[j+1].cost {
				wg.edges[j], wg.edges[j+1] = wg.edges[j+1], wg.edges[j]
			}
		}
	}
	
	parent := make([]int, wg.vertices)
	for i := 0; i < wg.vertices; i++ {
		parent[i] = -1
	}
	
	for i := 0; i < len(wg.edges); i++ {
		edge := wg.edges[i]
		
		x := wg.find(parent, edge.from)
		y := wg.find(parent, edge.to)
		
		if x != y {
			result = append(result, edge)
			<｜fim▁hole｜>
		}
	}
	
	return result
}

func (wg *WeightedGraph) Dijkstra(source int) []float64 {
	dist := make([]float64, wg.vertices)
	visited := make([]bool, wg.vertices)
	
	for i := 0; i < wg.vertices; i++ {
		dist[i] = math.Inf(1)
	}
	
	dist[source] = 0
	
	for count := 0; count < wg.vertices-1; count++ {
		u := -1
		minDist := math.Inf(1)
		
		for v := 0; v < wg.vertices; v++ {
			if !visited[v] && dist[v] < minDist {
				minDist = dist[v]
				u = v
			}
		}
		
		if u == -1 {
			break
		}
		
		visited[u] = true
		
		for _, edge := range wg.edges {
			if edge.from == u {
				v := edge.to
				if !visited[v] && dist[u]+edge.cost < dist[v] {
					dist[v] = dist[u] + edge.cost
				}
			} else if edge.to == u {
				v := edge.from
				if !visited[v] && dist[u]+edge.cost < dist[v] {
					dist[v] = dist[u] + edge.cost
				}
			}
		}
	}
	
	return dist
}

func main() {
	fmt.Println("Graph Algorithms Demo")
	fmt.Println("====================")
	
	// Unweighted Graph
	vertices := 6
	graph := NewGraph(vertices)
	
	edges := [][]int{
		{0, 1}, {0, 2}, {1, 3}, {1, 4},
		{2, 3}, {3, 4}, {3, 5}, {4, 5},
	}
	
	for _, edge := range edges {
		graph.AddEdge(edge[0], edge[1])
	}
	
	fmt.Println("\nUnweighted Graph:")
	fmt.Println("Vertices:", vertices)
	fmt.Println("Edges:", edges)
	
	fmt.Print("\nDFS traversal starting from vertex 0: ")
	visited := make([]bool, vertices)
	graph.DFS(0, visited)
	fmt.Println()
	
	fmt.Print("BFS traversal starting from vertex 0: ")
	graph.BFS(0)
	fmt.Println()
	
	// Weighted Graph
	fmt.Println("\nWeighted Graph:")
	weightedGraph := NewWeightedGraph(vertices)
	
	weightedEdges := [][]float64{
		{0, 1, 4}, {0, 2, 4}, {1, 2, 2}, {1, 3, 3},
		{1, 4, 2}, {2, 3, 1}, {2, 4, 4}, {3, 4, 3},
		{3, 5, 2}, {4, 5, 3},
	}
	
	for _, edge := range weightedEdges {
		weightedGraph.AddWeightedEdge(int(edge[0]), int(edge[1]), edge[2])
	}
	
	fmt.Println("Vertices:", vertices)
	fmt.Println("Edges (from, to, cost):", weightedEdges)
	
	// Kruskal's MST
	fmt.Println("\nKruskal's Minimum Spanning Tree:")
	mstKruskal := weightedGraph.KruskalMST()
	totalCost := 0.0
	for _, edge := range mstKruskal {
		fmt.Printf("  %d - %d: %.2f\n", edge.from, edge.to, edge.cost)
		totalCost += edge.cost
	}
	fmt.Printf("Total cost: %.2f\n", totalCost)
	
	// Dijkstra's Shortest Path
	fmt.Println("\nDijkstra's Shortest Path from vertex 0:")
	shortestPaths := weightedGraph.Dijkstra(0)
	for i := 0; i < vertices; i++ {
		if shortestPaths[i] == math.Inf(1) {
			fmt.Printf("  0 -> %d: Unreachable\n", i)
		} else {
			fmt.Printf("  0 -> %d: %.2f\n", i, shortestPaths[i])
		}
	}
}
