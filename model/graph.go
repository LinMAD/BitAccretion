package model

import "sync"

type (
	// VertexName represents name of vertex
	VertexName string

	// Graph of all infrastructure
	Graph struct {
		// TODO Refactor maps to slices cos Node will have name
		// Contains application vertex
		vertices map[VertexName]Node
		// Contains application vertex connection
		edges map[VertexName]map[VertexName]Node

		lock sync.RWMutex
	}

	// Node a container of value in graph
	Node struct {
		// TODO Make private
		Name         string
		Health       HealthState
		RequestCount int
		ErrorCount   int
		metadata     interface{}
	}
)

func (n Node) GetName() string {
	return n.Name
}

// NewGraph of with n nodes and edges
func NewGraph() *Graph {
	return &Graph{
		vertices: make(map[VertexName]Node),
		edges:    make(map[VertexName]map[VertexName]Node),
	}
}

// AddVertex by unique labeled vertex and return if it's added
func (g *Graph) AddVertex(name VertexName, node Node) bool {
	g.lock.Lock()
	defer g.lock.Unlock()

	if _, ok := g.vertices[name]; ok {
		return false
	}

	g.vertices[name] = node

	return true
}

// AddEdge between vertex 'from' and vertex 'to'
func (g *Graph) AddEdge(from, to VertexName, node Node) {
	g.AddVertex(from, node)
	g.AddVertex(to, node)

	g.lock.Lock()
	defer g.lock.Unlock()

	if _, ok := g.edges[from]; !ok {
		g.edges[from] = make(map[VertexName]Node)
	}

	g.edges[from][to] = node
}

// GetAllVerticesLabels returns all registered vertices in graph
func (g *Graph) GetAllVerticesLabels() (names []VertexName) {
	names = make([]VertexName, len(g.vertices))

	g.lock.RLock()
	defer g.lock.RUnlock()

	var i int
	for n := range g.vertices {
		names[i] = n
		i++
	}

	return names
}

// GetVertexEdges returns all related edges
func (g *Graph) GetVertexEdges(vl VertexName) map[VertexName]Node {
	return g.edges[vl]
}

// GetVertex returns vertex data
func (g *Graph) GetVertex(vl VertexName) Node {
	return g.vertices[vl]
}
