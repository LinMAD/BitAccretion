package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/LinMAD/BitAccretion/logger"
	"github.com/LinMAD/BitAccretion/model"
	"github.com/LinMAD/BitAccretion/plugin/newrelic/worker"
	"github.com/LinMAD/BitAccretion/provider"
)

// NRConfig addition for config required for New Relic
type NRConfig struct {
	// APIKey of NewRelic
	APIKey string `json:"api_key"`
	// APPSets to survey in NewRelic
	APPSets []APPSet `json:"app_sets"`
}

// APPSet registered in New Relic
type APPSet struct {
	AppDetails
}

// AppDetails detailed information about application in New Relic
type AppDetails struct {
	// Name of application in NewRelic system (case sensitive)
	Name string `json:"name"`
	// ID of application in NewRelic
	ID string `json:"id"`
	// RelicMetrics to get from API for surveyed system
	RelicMetrics []string `json:"relic_metrics"`
}


// ProviderNewRelic base structure for plugin workflow
type ProviderNewRelic struct {
	// Config for new relic plugin
	Config NRConfig
	// log messages in runtime
	log logger.ILogger
	// worker to harvesting API for metric data
	worker *worker.RelicWorker
	// pluginHealth current state of plugin
	pluginHealth model.HealthState
	// systemGraph is observable data
	systemGraph model.Graph
}

// LoadConfig of new relic plugin
func (nr *ProviderNewRelic) LoadConfig(pathToConfig string) error {
	configFile, openErr := os.Open(pathToConfig)
	if openErr != nil {
		return openErr
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)

	var relicConf NRConfig
	decodeErr := jsonParser.Decode(&relicConf)
	if decodeErr != nil {
		return decodeErr
	}

	nr.Config = relicConf

	return nil
}

// Boot setups all configuration and dependencies
func (nr *ProviderNewRelic) Boot(l logger.ILogger) error {
	relicWorker, relicWorkerErr := worker.NewRelicWorker(nr.Config.APIKey)
	if relicWorkerErr != nil {
		return relicWorkerErr
	}

	nr.worker = relicWorker
	nr.log = l

	return nil
}

// DispatchMonitoredData returns latest assembled graph
func (nr *ProviderNewRelic) DispatchMonitoredData() (model.Graph, error) {
	g := nr.systemGraph

	appList := g.GetAllVertices()
	appCount := int8(len(appList))

	var wg sync.WaitGroup
	var w int8

	for w = 0; w < appCount; w++ {
		wg.Add(1)

		go func(w int8) {
			defer wg.Done()

			app := appList[w]
			nr.log.Normal(fmt.Sprintf("Surveying metrics of '%s'", app.Name))
			if app.MetaData == nil {
				return
			}

			// Get base application data
			fetchedMetrics := nr.worker.CollectApplicationHostMetrics(
				app.MetaData.(AppDetails).ID,
				app.MetaData.(AppDetails).RelicMetrics,
			)

			for _, host := range fetchedMetrics.HostMetrics {
				app.Metric.RequestCount += host.Metrics.RequestCount
				app.Metric.ErrorCount += host.Metrics.ErrorCount
			}
		}(w)
	}

	wg.Wait()

	return g, nil
}

// ProvideHealth will return health of plugin if New Relic API alive\reachable
func (nr *ProviderNewRelic) ProvideHealth() model.HealthState {
	return nr.pluginHealth
}

// createBaseGraph from given configuration
func (nr *ProviderNewRelic) createBaseGraph() {
	nr.systemGraph = *model.NewGraph()

	for _, set := range nr.Config.APPSets {
		nr.systemGraph.AddVertex(
			model.VertexName(set.Name),
			model.Node{
				Name:   set.Name,
				Health: model.HealthNormal,
				Metric: model.SystemMetric{},
				MetaData: AppDetails{
					Name:         set.Name,
					ID:           set.ID,
					RelicMetrics: set.RelicMetrics,
				},
			},
		)
	}
}

// NewProvider returns instance with implemented interface
func NewProvider() provider.IProvider {
	return &ProviderNewRelic{pluginHealth: model.HealthNormal}
}
