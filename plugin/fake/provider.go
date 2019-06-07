package main

import (
	"math/rand"

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

// DispatchMonitoredData with dummy data
func (f *FakeProvider) DispatchMonitoredData() (model.Graph, error) {
	return GetStubGraph(), nil
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
func GetStubNodes() []*model.Node {
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
			Metric: model.SystemMetric{
				RequestCount: float32(rand.Intn(1000)),
				ErrorCount:   float32(rand.Intn(1000)),
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
