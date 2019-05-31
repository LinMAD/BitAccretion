package model

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type GraphTestSuite struct {
	suite.Suite
	graph *Graph
}

func (t *GraphTestSuite) SetupTest() {
	t.graph = NewGraph()
}

func (t *GraphTestSuite) TearDownAllSuite() {
	t.graph = nil
}

func TestRunGraphTestSuite(t *testing.T) {
	suite.Run(t, new(GraphTestSuite))
}

func (t *GraphTestSuite) TestVertexEdgeRelation() {
	var edges map[VertexName]Node
	testNode := Node{}

	t.graph.AddEdge("beta", "theta", testNode)
	t.graph.AddEdge("theta", "kilo", testNode)

	edges = t.graph.GetVertexEdges("beta")
	if _, ok := edges["theta"]; !ok {
		t.Fail("Edge not found from `beta` to `theta` vertex")
	}

	edges = t.graph.GetVertexEdges("theta")
	if _, ok := edges["kilo"]; !ok {
		t.Fail("Edge not found from `theta` to `kilo` vertex")
	}
}

func (t *GraphTestSuite) TestAddVertex() {
	t.graph.AddVertex(VertexName("hotel"), Node{})
	vertexes := t.graph.GetAllVerticesLabels()

	assert.Equal(t.Suite.T(), len(vertexes), 1)
	assert.Equal(t.Suite.T(), vertexes[0], VertexName("hotel"))
	assert.NotEqual(t.Suite.T(), vertexes[0], "hotel")

	assert.False(t.Suite.T(), t.graph.AddVertex(VertexName("hotel"), Node{}))
}

type testAppGraph struct {
	Name      string
	NestedApp []testAppGraph
}

func graphProvider(g *Graph) {
	layer3 := make([]testAppGraph, 2)
	layer3[0] = testAppGraph{
		Name: "Echo",
	}
	layer3[1] = testAppGraph{
		Name: "Foxtrot",
	}
	layer2 := make([]testAppGraph, 1)
	layer2[0] = testAppGraph{
		Name:      "Delta",
		NestedApp: layer3,
	}

	// Set first layer
	layer := make([]testAppGraph, 2)
	layer[0] = testAppGraph{
		Name: "Bravo",
	}
	layer[1] = testAppGraph{
		Name:      "Charlie",
		NestedApp: layer2,
	}

	// Combine tree
	infraConfig := testAppGraph{
		Name:      "Alfa",
		NestedApp: layer,
	}

	infraLoaderHelper(infraConfig, g)
}

func (t *GraphTestSuite) TestGraphStructure() {
	appNamesList := map[string]bool{
		"Alfa": false, "Bravo": false,
		"Charlie": false, "Delta": false,
		"Echo": false, "Foxtrot": false,
	}

	graphProvider(t.graph)

	// Validate
	vertices := t.graph.GetAllVerticesLabels()
	for _, v := range vertices {
		if _, ok := appNamesList[string(v)]; !ok {
			assert.Failf(t.T(), "Expected to be found vertex name", string(v))
		}
	}

	charlieEdge := t.graph.GetVertexEdges("Charlie")
	if _, ok := charlieEdge["Delta"]; !ok {
		assert.Failf(t.T(), "Expected to be found next edge from `Charlie` to `Delta`", "Edge `Delta` not found")
	}

	deltaEdge := t.graph.GetVertexEdges("Delta")
	if _, ok := deltaEdge["Echo"]; !ok {
		assert.Failf(t.T(), "Expected to be found next edge from `Delta` to `Echo`", "Edge `Echo` not found")
	}
	if _, ok := deltaEdge["Foxtrot"]; !ok {
		assert.Failf(t.T(), "Expected to be found next edge from `Delta` to `Foxtrot`", "Edge `Foxtrot` not found")
	}
}

func (t *GraphTestSuite) TestGetAllVertices() {
	graphProvider(t.graph)
	vertices := t.graph.GetAllVertices()

	assert.NotZero(t.T(), vertices, )
}

func infraLoaderHelper(infra testAppGraph, graph *Graph) {
	if infra.NestedApp == nil {
		graph.AddVertex(VertexName(infra.Name), Node{})
		return
	}

	for _, nested := range infra.NestedApp {
		graph.AddEdge(VertexName(infra.Name), VertexName(nested.Name), Node{})

		infraLoaderHelper(nested, graph)
	}
}
