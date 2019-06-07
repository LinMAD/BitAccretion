package main

import (
	"math/rand"

	"github.com/LinMAD/BitAccretion/logger"
	"github.com/LinMAD/BitAccretion/model"
	"github.com/LinMAD/BitAccretion/provider"
)

// FakeProvider randomly generates dummy data
type FakeProvider struct {
	pluginHealth model.HealthState
}

// LoadConfig stub
func (f *FakeProvider) LoadConfig(pathToConfig string) error {
	return nil
}

// Boot stub
func (f *FakeProvider) Boot() error {
	return nil
}

// DispatchGraph prepared graph of monitored systems
func (f *FakeProvider) DispatchGraph() (model.Graph, error) {
	return GetStubGraph(false), nil
}

// FetchNewData with dummy data
func (f *FakeProvider) FetchNewData(log logger.ILogger) (model.Graph, error) {
	log.Normal("Generating fake data...")

	return GetStubGraph(true), nil
}

// ProvideHealth immortal
func (f *FakeProvider) ProvideHealth() model.HealthState {
	return f.pluginHealth
}

// NewProvider implementation
func NewProvider() provider.IProvider {
	return &FakeProvider{pluginHealth: model.HealthNormal}
}

// GetStubNodes generated dummy nodes with data
func GetStubNodes(isRandom bool) []*model.Node {
	sysNames := []string{
		"Lipstick", "Steward",
		"Siege Engine", "Homesick",
		"Gray Knife", "Jungle Paladin",
		"Urban Scorpion", "Magnet",
		"Blockade", "Boomstick",
		"Orange Jack", "Red Winter",
	}

	nodes := make([]*model.Node, len(sysNames))
	for i := 0; i < len(sysNames); i++ {
		nodes[i] = &model.Node{
			Name:   sysNames[i],
			Health: getRandomHealthState(),
		}

		if isRandom {
			nodes[i].Metric = model.SystemMetric{
				RequestCount: float32(rand.Intn(10000)),
				ErrorCount:   float32(rand.Intn(50)),
			}
		}
	}

	return nodes
}

// GetStubGraph generated dummy graph with nodes
func GetStubGraph(isRandom bool) model.Graph {
	g := model.NewGraph()

	for _, n := range GetStubNodes(isRandom) {
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
