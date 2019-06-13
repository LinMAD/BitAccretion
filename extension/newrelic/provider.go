package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/LinMAD/BitAccretion/extension"
	"github.com/LinMAD/BitAccretion/extension/newrelic/worker"
	"github.com/LinMAD/BitAccretion/logger"
	"github.com/LinMAD/BitAccretion/model"
	"github.com/LinMAD/BitAccretion/util"
)

// NRConfig addition for config required for New Relic
type NRConfig struct {
	// APIKey of NewRelic
	APIKey string `json:"api_key"`
	// APPSets to survey in NewRelic
	APPSets []APPSet `json:"app_sets"`
	// HealthSensitivity conversion to mark health of vertex
	HealthSensitivity model.HealthSensitivity `json:"health_sensitivity"`
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

// ProviderNewRelic base structure for extension workflow
type ProviderNewRelic struct {
	// Config for new relic extension
	Config NRConfig
	// worker to harvesting API for metric data
	worker *worker.RelicWorker
	// pluginHealth current state of extension
	pluginHealth model.HealthState
}

// LoadConfig of new relic extension
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
func (nr *ProviderNewRelic) Boot() error {
	relicWorker, relicWorkerErr := worker.NewRelicWorker(nr.Config.APIKey)
	if relicWorkerErr != nil {
		return relicWorkerErr
	}

	nr.worker = relicWorker

	return nil
}

// DispatchGraph prepared new relic graph of monitored systems
func (nr *ProviderNewRelic) DispatchGraph() (model.Graph, error) {
	return *nr.prepareGraph(), nil
}

// FetchNewData returns latest assembled graph
func (nr *ProviderNewRelic) FetchNewData(log logger.ILogger) (model.Graph, error) {
	g := nr.prepareGraph()

	appList := g.GetAllVertices()
	appCount := int8(len(appList))

	log.Debug(fmt.Sprintf("Harvesting data from NewRelic API..."))

	var wg sync.WaitGroup
	var w int8
	for w = 0; w < appCount; w++ {
		wg.Add(1)

		go func(w int8) {
			defer wg.Done()

			app := appList[w]
			if app.MetaData == nil {
				return
			}

			// Get base application data
			fetchedMetrics := nr.worker.CollectApplicationHostMetrics(
				log,
				app.MetaData.(AppDetails).ID,
				app.MetaData.(AppDetails).RelicMetrics,
			)

			app.Metric = model.SystemMetric{}
			for _, host := range fetchedMetrics.HostMetrics {
				app.Metric.RequestCount += host.Metrics.RequestCount
				app.Metric.ErrorCount += host.Metrics.ErrorCount
			}

			app.Health = util.GetMetricHealthByValue(&app.Metric, &nr.Config.HealthSensitivity)
		}(w)
	}

	wg.Wait()

	return *g, nil
}

// ProvideHealth will return health of extension if New Relic API alive\reachable
func (nr *ProviderNewRelic) ProvideHealth() model.HealthState {
	return nr.pluginHealth
}

// prepareGraph from given configuration
func (nr *ProviderNewRelic) prepareGraph() (g *model.Graph) {
	g = model.NewGraph()

	for _, set := range nr.Config.APPSets {
		g.AddVertex(
			model.VertexName(set.Name),
			&model.Node{
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

	return
}

// NewProvider returns instance with implemented interface
func NewProvider() extension.IProvider {
	return &ProviderNewRelic{pluginHealth: model.HealthNormal}
}
