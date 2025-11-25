package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sort"
	"sync"
)

// 简化版图计算系统实现

// Vertex 顶点
type Vertex struct {
	ID    string  `json:"id"`
	Value float64 `json:"value"`
}

// Edge 边
type Edge struct {
	ID     string  `json:"id"`
	Source string  `json:"source"`
	Target string  `json:"target"`
	Weight float64 `json:"weight"`
}

// Graph 图
type Graph struct {
	Vertices map[string]*Vertex  `json:"vertices"`
	Edges    map[string]*Edge    `json:"edges"`
	InEdges  map[string][]string `json:"in_edges"`
	OutEdges map[string][]string `json:"out_edges"`
	mu       sync.RWMutex
}

// NewGraph 创建新图
func NewGraph() *Graph {
	return &Graph{
		Vertices: make(map[string]*Vertex),
		Edges:    make(map[string]*Edge),
		InEdges:  make(map[string][]string),
		OutEdges: make(map[string][]string),
	}
}

// AddVertex 添加顶点
func (g *Graph) AddVertex(vertex *Vertex) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.Vertices[vertex.ID] = vertex
	if _, exists := g.InEdges[vertex.ID]; !exists {
		g.InEdges[vertex.ID] = make([]string, 0)
	}
	if _, exists := g.OutEdges[vertex.ID]; !exists {
		g.OutEdges[vertex.ID] = make([]string, 0)
	}
}

// AddEdge 添加边
func (g *Graph) AddEdge(edge *Edge) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.Edges[edge.ID] = edge
	g.OutEdges[edge.Source] = append(g.OutEdges[edge.Source], edge.ID)
	g.InEdges[edge.Target] = append(g.InEdges[edge.Target], edge.ID)
}

// GetVertex 获取顶点
func (g *Graph) GetVertex(id string) (*Vertex, bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	vertex, exists := g.Vertices[id]
	return vertex, exists
}

// GetVertices 获取所有顶点
func (g *Graph) GetVertices() []*Vertex {
	g.mu.RLock()
	defer g.mu.RUnlock()

	vertices := make([]*Vertex, 0, len(g.Vertices))
	for _, vertex := range g.Vertices {
		vertices = append(vertices, vertex)
	}
	return vertices
}

// GetInEdges 获取顶点的入边
func (g *Graph) GetInEdges(vertexID string) []*Edge {
	g.mu.RLock()
	defer g.mu.RUnlock()

	edgeIDs := g.InEdges[vertexID]
	edges := make([]*Edge, 0, len(edgeIDs))
	for _, edgeID := range edgeIDs {
		if edge, exists := g.Edges[edgeID]; exists {
			edges = append(edges, edge)
		}
	}
	return edges
}

// GetOutEdges 获取顶点的出边
func (g *Graph) GetOutEdges(vertexID string) []*Edge {
	g.mu.RLock()
	defer g.mu.RUnlock()

	edgeIDs := g.OutEdges[vertexID]
	edges := make([]*Edge, 0, len(edgeIDs))
	for _, edgeID := range edgeIDs {
		if edge, exists := g.Edges[edgeID]; exists {
			edges = append(edges, edge)
		}
	}
	return edges
}

// PageRank PageRank算法
type PageRank struct {
	dampingFactor float64
	iterations    int
	tolerance     float64
}

// NewPageRank 创建新的PageRank算法
func NewPageRank(dampingFactor float64, iterations int, tolerance float64) *PageRank {
	return &PageRank{
		dampingFactor: dampingFactor,
		iterations:    iterations,
		tolerance:     tolerance,
	}
}

// Execute 执行PageRank算法
func (pr *PageRank) Execute(graph *Graph, params map[string]interface{}) map[string]interface{} {
	vertices := graph.GetVertices()

	// 初始化PageRank值
	prValues := make(map[string]float64)
	for _, vertex := range vertices {
		prValues[vertex.ID] = 1.0 / float64(len(vertices))
	}

	// 迭代计算
	for i := 0; i < pr.iterations; i++ {
		newPRValues := make(map[string]float64)

		// 计算每个顶点的新PageRank值
		for _, vertex := range vertices {
			sum := 0.0

			// 收集所有指向该顶点的边的贡献
			inEdges := graph.GetInEdges(vertex.ID)
			for _, edge := range inEdges {
				sourceVertex, exists := graph.GetVertex(edge.Source)
				if exists {
					sourceOutEdges := graph.GetOutEdges(sourceVertex.ID)
					sourceOutDegree := len(sourceOutEdges)
					if sourceOutDegree > 0 {
						sum += prValues[sourceVertex.ID] / float64(sourceOutDegree)
					}
				}
			}

			// 计算新的PageRank值
			newPRValues[vertex.ID] = (1.0-pr.dampingFactor)/float64(len(vertices)) + pr.dampingFactor*sum
		}

		// 检查收敛
		maxDiff := 0.0
		for id, value := range prValues {
			diff := math.Abs(value - newPRValues[id])
			if diff > maxDiff {
				maxDiff = diff
			}
		}

		// 更新PageRank值
		prValues = newPRValues

		// 如果收敛则提前退出
		if maxDiff < pr.tolerance {
			break
		}
	}

	// 更新顶点值
	for _, vertex := range vertices {
		vertex.Value = prValues[vertex.ID]
	}

	// 排序顶点
	sortedVertices := make([]*Vertex, len(vertices))
	copy(sortedVertices, vertices)
	sort.Slice(sortedVertices, func(i, j int) bool {
		return sortedVertices[i].Value > sortedVertices[j].Value
	})

	// 准备结果
	result := make(map[string]interface{})
	result["iterations"] = pr.iterations

	topVertices := make([]map[string]interface{}, 0, 5)
	for i := 0; i < 5 && i < len(sortedVertices); i++ {
		vertex := sortedVertices[i]
		topVertex := map[string]interface{}{
			"id":    vertex.ID,
			"value": vertex.Value,
		}
		topVertices = append(topVertices, topVertex)
	}
	result["top_vertices"] = topVertices

	return result
}

// Name 获取算法名称
func (pr *PageRank) Name() string {
	return "PageRank"
}

// SimpleGraphEngine 简化版图计算引擎
type SimpleGraphEngine struct {
	graph *Graph
}

// NewSimpleGraphEngine 创建新的简化版图计算引擎
func NewSimpleGraphEngine() *SimpleGraphEngine {
	return &SimpleGraphEngine{
		graph: NewGraph(),
	}
}

// AddVertex 添加顶点
func (e *SimpleGraphEngine) AddVertex(vertex *Vertex) {
	e.graph.AddVertex(vertex)
}

// AddEdge 添加边
func (e *SimpleGraphEngine) AddEdge(edge *Edge) {
	e.graph.AddEdge(edge)
}

// ExecuteAlgorithm 执行算法
func (e *SimpleGraphEngine) ExecuteAlgorithm(algorithm interface{}, params map[string]interface{}) map[string]interface{} {
	if pr, ok := algorithm.(*PageRank); ok {
		return pr.Execute(e.graph, params)
	}

	return map[string]interface{}{
		"error": "不支持的算法",
	}
}

// GetStats 获取统计信息
func (e *SimpleGraphEngine) GetStats() map[string]interface{} {
	stats := make(map[string]interface{})
	stats["vertex_count"] = len(e.graph.Vertices)
	stats["edge_count"] = len(e.graph.Edges)<｜fim▁hole｜>

	return stats
}

func main() {
	// 创建简化版图计算引擎
	engine := NewSimpleGraphEngine()

	// 添加一些顶点和边
	for i := 1; i <= 10; i++ {
		vertex := &Vertex{
			ID:    fmt.Sprintf("v%d", i),
			Value: rand.Float64(),
		}
		engine.AddVertex(vertex)
	}

	// 添加一些边
	for i := 1; i <= 10; i++ {
		for j := 1; j <= 2; j++ {
			target := (i + j) % 10
			if target == 0 {
				target = 10
			}

			edge := &Edge{
				ID:     fmt.Sprintf("e%d_%d", i, target),
				Source: fmt.Sprintf("v%d", i),
				Target: fmt.Sprintf("v%d", target),
				Weight: rand.Float64(),
			}
			engine.AddEdge(edge)
		}
	}

	// 执行PageRank算法
	pagerank := NewPageRank(0.85, 10, 0.0001)
	pagerankResult := engine.ExecuteAlgorithm(pagerank, map[string]interface{}{})
	log.Printf("PageRank结果: %+v", pagerankResult)

	// 获取统计信息
	stats := engine.GetStats()
	log.Printf("引擎统计: %+v", stats)

	fmt.Println("简化版图计算系统示例完成")
}
