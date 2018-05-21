package graph

import (
	"github.com/LinMAD/BitAccretion/core/assembly/structure"
	"time"
)

// VertexLabel represents name of vertex
type VertexLabel string

// Graph of all infrastructure
type Graph struct {
	// Contains application vertex
	vertices map[VertexLabel]*structure.VNode
	// Contains application vertex connection
	edges map[VertexLabel]map[VertexLabel]*structure.VNodeConnection
}

// NewGraph of with n nodes and edges
func NewGraph() *Graph {
	return &Graph{
		vertices: make(map[VertexLabel]*structure.VNode),
		edges:    make(map[VertexLabel]map[VertexLabel]*structure.VNodeConnection),
	}
}

// AddVertex by unique labeled vertex and return if it's added
func (g *Graph) AddVertex(name VertexLabel) bool {
	if _, ok := g.vertices[name]; ok {
		return false
	}

	g.vertices[name] = &structure.VNode{
		Name:        string(name),
		DisplayName: string(name),
		Updated:     time.Now().UnixNano(),
	}

	return true
}

// AddEdge between vertex 'from' and vertex 'to'
func (g *Graph) AddEdge(from, to VertexLabel) {
	if _, ok := g.vertices[from]; !ok {
		g.AddVertex(from)
	}

	if _, ok := g.vertices[to]; !ok {
		g.AddVertex(to)
	}

	if _, ok := g.edges[from]; !ok {
		g.edges[from] = make(map[VertexLabel]*structure.VNodeConnection)
	}

	g.edges[from][to] = &structure.VNodeConnection{
		Source: string(from),
		Target: string(to),
	}
}

// GetAllVertices returns all registered vertices in graph
func (g *Graph) GetAllVertices() (names []VertexLabel) {
	for n := range g.vertices {
		names = append(names, n)
	}

	return names
}

// GetVertexEdges returns all related edges
func (g *Graph) GetVertexEdges(name VertexLabel) map[VertexLabel]*structure.VNodeConnection {
	return g.edges[name]
}

// GetVertex returns vertex data
func (g *Graph) GetVertex(name VertexLabel) *structure.VNode {
	return g.vertices[name]
}
