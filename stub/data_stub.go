package stub

import (
	"math/rand"

	"github.com/LinMAD/BitAccretion/model"
)

// GetStubNodes generated dummy nodes with data
func GetStubNodes() []model.Node {
	sysNames := []string{
		"Lipstick", "Steward",
		"Siege Engine", "Homesick",
		"Gray Knife", "Jungle Paladin",
		"Urban Scorpion", "Magnet",
		"Blockade", "Boomstick",
		"Orange Jack", "Red Winter",
	}

	nodes := make([]model.Node, len(sysNames))
	for i := 0; i < len(sysNames); i++ {
		nodes[i] = model.Node{
			Name:   sysNames[i],
			Health: getRandomHealthState(),
			Metric: model.SystemMetric{
				RequestCount: rand.Intn(1000),
				ErrorCount:   rand.Intn(1000),
			},
		}
	}

	return nodes
}

// GetStubGraph generated dummy graph with nodes
func GetStubGraph() model.Graph {
	g := model.NewGraph()

	for _, n := range GetStubNodes() {
		g.AddVertex(model.VertexName(n.Name), n)
	}

	return *g
}

func getRandomHealthState() model.HealthState {
	any := model.HealthState(rand.Intn(len(model.HealthStatesMap)))

	for h := range model.HealthStatesMap {
		if any == h {
			return h
		}
	}

	return model.HealthNormal
}
