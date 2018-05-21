package assembly

import (
	"encoding/json"
	"github.com/LinMAD/BitAccretion/core/assembly/graph"
	"github.com/LinMAD/BitAccretion/core/assembly/structure"
	"io/ioutil"
	"log"
)

// INTERNET node name in graph (general)
const INTERNET = "INTERNET"

// InfraObject represents infrastructure object with relation
type InfraObject struct {
	// Name of object
	Name string
	// NestedInfraObject related object to current
	NestedInfraObject []InfraObject
	// Details some additional data related to object
	Details interface{}
}

// HealthSensitivity represents metrics sensitivity determination
type HealthSensitivity struct {
	Danger  float32 `json:"danger"`
	Warning float32 `json:"warning"`
}

// MakeInfrastructureGraph creates graph based on InfraObject
func MakeInfrastructureGraph(infra []InfraObject) *graph.Graph {
	g := graph.NewGraph()

	for _, object := range infra {
		buildGraph(object, g)
	}

	return g
}

// ConvertToVizceral converts graph to json representation of vizceral structure
func ConvertToVizceral(g *graph.Graph, healthSens HealthSensitivity) structure.VRegionGraph {
	regeon := structure.VRegionGraph{
		Renderer: structure.VGlobalRenderer,
		Name:     "edge", // vizceral global name
	}

	inetNode := structure.VNode{
		Renderer:    structure.VRegionRenderer,
		Name:        INTERNET,
		DisplayName: INTERNET,
		Metadata: structure.VMeta{
			Streaming: 1,
		},
	}

	for _, name := range g.GetAllVertices() {
		// Collect total metrics
		var internetMetrics structure.VMetricLevels
		edges := g.GetVertexEdges(name)
		for _, edge := range edges {
			// Add relation for vizceral
			sourceVertex := *g.GetVertex(graph.VertexLabel(edge.Source))
			sourceVertex.Renderer = structure.VRegionRenderer

			// Add nested scaled nodes
			for _, relatedEdges := range edges {
				internetMetrics.Normal += relatedEdges.Metrics.Normal
				internetMetrics.Warning += relatedEdges.Metrics.Warning
				internetMetrics.Danger += relatedEdges.Metrics.Danger

				sourceConnect := structure.VNodeConnection{
					Source:  relatedEdges.Source,
					Target:  relatedEdges.Target,
					Metrics: relatedEdges.Metrics,
					Class:   GetMetricHealth(healthSens, relatedEdges.Metrics),
				}

				nestedVertex := *g.GetVertex(graph.VertexLabel(relatedEdges.Target))
				nestedVertex.Class = GetMetricHealth(healthSens, relatedEdges.Metrics)
				sourceVertex.Nodes = append(sourceVertex.Nodes, nestedVertex)
				sourceVertex.Connections = append(sourceVertex.Connections, sourceConnect)
			}

			// Set internet connection
			internetConn := structure.VNodeConnection{
				Source:  INTERNET,
				Target:  edge.Source,
				Metrics: internetMetrics,
				Class:   GetMetricHealth(healthSens, internetMetrics),
			}

			if isNeedToNotice, notice := GetHealthNotice(internetConn.Class); isNeedToNotice {
				internetConn.Notices = append(internetConn.Notices, notice)
			}

			// Set source of nested node inside for vizceral
			sourceNodeNested := structure.VNode{
				Renderer:    structure.VRegionRenderer,
				Name:        edge.Source,
				DisplayName: edge.Source,
				Class:       internetConn.Class,
			}

			sourceVertex.Nodes = append(sourceVertex.Nodes, sourceNodeNested)
			sourceVertex.Class = internetConn.Class

			regeon.Nodes = append(regeon.Nodes, sourceVertex)
			regeon.Connections = append(regeon.Connections, internetConn)

			break
		}
	}

	regeon.Nodes = append(regeon.Nodes, inetNode)

	return regeon
}

// buildGraph recursively builds graph for given infrastructure object structure
func buildGraph(object InfraObject, infraGraph *graph.Graph) {
	var vertex *structure.VNode
	if object.NestedInfraObject == nil {
		// Create vertex
		infraGraph.AddVertex(graph.VertexLabel(object.Name))

		vertex := infraGraph.GetVertex(graph.VertexLabel(object.Name))
		vertex.SystemDetails = object.Details

		return
	}

	for _, nested := range object.NestedInfraObject {
		// Relate 2 objects
		infraGraph.AddEdge(graph.VertexLabel(object.Name), graph.VertexLabel(nested.Name))

		// Add vertices details
		vertex = infraGraph.GetVertex(graph.VertexLabel(object.Name))
		vertex.SystemDetails = object.Details
		vertex = infraGraph.GetVertex(graph.VertexLabel(nested.Name))
		vertex.SystemDetails = object.Details

		// Move to nested object
		buildGraph(nested, infraGraph)
	}
}

// GetMetricHealth return node health status by given metrics
func GetMetricHealth(sensitivity HealthSensitivity, appMetrics structure.VMetricLevels) structure.VClassType {
	// Define node health by metrics
	if appMetrics.Danger > sensitivity.Danger {
		return structure.VDanger
	} else if appMetrics.Warning > sensitivity.Warning || appMetrics.Danger > 0 {
		return structure.VWarning
	} else {
		return structure.VNormal
	}
}

// GetHealthNotice will return notice structure if node has bad health as error conversion
func GetHealthNotice(class structure.VClassType) (isWarn bool, notice structure.VNotice) {
	if class == structure.VDanger {
		notice.Title = "Error conversion high"
		notice.Severity = 2

		isWarn = true
	} else if class == structure.VWarning {
		notice.Title = "Error conversion medium"
		notice.Severity = 1

		isWarn = true
	}

	return
}

// WriteToJSON convert data to json format
func WriteToJSON(data interface{}) []byte {
	vJSON, err := json.Marshal(data)
	if err != nil {
		log.Print(err.Error())
	}

	return vJSON
}

// WriteJSONToFile write data to given file in json format
func WriteJSONToFile(pathFile string, data interface{}) {
	if data == nil {
		return
	}

	err := ioutil.WriteFile(pathFile, []byte(WriteToJSON(data)), 0644)
	if err != nil {
		log.Fatalf("Unable to write to file: %s -> %v", pathFile, err.Error())
	}
}
