package assembly

import (
	"encoding/json"
	"github.com/LinMAD/BitAccretion/core/assembly/graph"
	"github.com/LinMAD/BitAccretion/core/assembly/structure"
	"github.com/stretchr/testify/suite"
	"testing"
)

type AssemblyTestSuite struct {
	suite.Suite
	rawStructure2D   map[string]string
	infrastructure2D []InfraObject
}

func (t *AssemblyTestSuite) SetupTest() {
	t.rawStructure2D = map[string]string{
		"Alfa":    "Bravo",
		"Charlie": "Delta",
		"Echo":    "Foxtrot",
		"Golf":    "Hotel",
	}

	for name, nested := range t.rawStructure2D {
		nestedApp := make([]InfraObject, 1)
		nestedApp[0] = InfraObject{Name: nested}

		infApp := InfraObject{
			Name:              name,
			NestedInfraObject: nestedApp,
		}

		t.infrastructure2D = append(t.infrastructure2D, infApp)
	}

}

func TestRunGraphTestSuite(t *testing.T) {
	suite.Run(t, new(AssemblyTestSuite))
}

func (t *AssemblyTestSuite) TestCreateInfrastructureGraph2D() {
	infraGraph := MakeInfrastructureGraph(t.infrastructure2D)

	for name, nested := range t.rawStructure2D {
		edge := infraGraph.GetVertexEdges(graph.VertexLabel(name))
		if _, ok := edge[graph.VertexLabel(nested)]; !ok {
			t.Failf("Edge (%s) not found for vertex (%s)", name, nested)
			return
		}
	}
}

func (t *AssemblyTestSuite) TestCreateInfrastructureGraph3D() {
	infra := make([]InfraObject, 1)
	infraLayer2 := make([]InfraObject, 1)
	infraLayer3 := make([]InfraObject, 1)

	infraLayer3[0] = InfraObject{Name: "Charlie"}
	infraLayer2[0] = InfraObject{
		Name:              "Bravo",
		NestedInfraObject: infraLayer3,
	}
	infra[0] = InfraObject{
		Name:              "Alfa",
		NestedInfraObject: infraLayer2,
	}

	infraGraph := MakeInfrastructureGraph(infra)

	edge1 := infraGraph.GetVertexEdges(graph.VertexLabel("Alfa"))
	if _, ok := edge1[graph.VertexLabel("Bravo")]; !ok {
		t.Failf("Edge (%s) not found for vertex (%s)", "Alfa", "Bravo")
		return
	}

	edge2 := infraGraph.GetVertexEdges(graph.VertexLabel("Bravo"))
	if _, ok := edge2[graph.VertexLabel("Charlie")]; !ok {
		t.Failf("Edge (%s) not found for vertex (%s)", "Bravo", "Charlie")
	}
}

func (t *AssemblyTestSuite) TestMetricsHealth() {
	healthSens := HealthSensitivity{
		Danger:  10,
		Warning: 5,
	}

	metrics := structure.VMetricLevels{
		Danger:  42,
		Warning: 2,
		Normal:  1,
	}

	health := GetMetricHealth(healthSens, metrics)
	t.Assert().Equal(structure.VDanger, health)

	metrics = structure.VMetricLevels{
		Danger:  9,
		Warning: 42,
		Normal:  1,
	}

	health = GetMetricHealth(healthSens, metrics)
	t.Assert().Equal(structure.VWarning, health)

	metrics = structure.VMetricLevels{
		Danger:  9,
		Warning: 0,
		Normal:  1,
	}

	health = GetMetricHealth(healthSens, metrics)
	t.Assert().Equal(structure.VWarning, health)

	metrics = structure.VMetricLevels{
		Danger:  0,
		Warning: 0,
		Normal:  1,
	}

	health = GetMetricHealth(healthSens, metrics)
	t.Assert().Equal(structure.VNormal, health)
}

type testData struct {
	Name string
	Val  int
}

func (t *AssemblyTestSuite) TestWriteToJSON() {
	data := testData{
		Name: "Unit",
		Val:  42,
	}

	jByte := WriteToJSON(data)

	var decoded testData
	err := json.Unmarshal(jByte, &decoded)
	if err != nil {
		t.Failf("Unable to decode encoded json", "err: %v", err)
		return
	}

	t.Assert().Equal("Unit", decoded.Name)
	t.Assert().Equal(42, decoded.Val)
}
