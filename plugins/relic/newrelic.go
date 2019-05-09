package main

import (
	"encoding/json"
	"github.com/LinMAD/BitAccretion/core"
	"github.com/LinMAD/BitAccretion/core/assembly"
	"github.com/LinMAD/BitAccretion/core/assembly/graph"
	"github.com/LinMAD/BitAccretion/core/assembly/structure"
	"github.com/LinMAD/BitAccretion/plugins/relic/worker"
	"log"
	"os"
	"plugin"
	"sync"
	"time"
)

const (
	nrTag = "PRC_NEW_RELIC"
)

// NewRelicProcessor base structure for plugin workflow
type NewRelicProcessor struct {
	Config       NRConfig
	graph        *graph.Graph
	vRegionGraph structure.VRegionGraph
	relicWorker  *worker.RelicWorker
	sound        safeSoundAlert
}

// safeSoundAlert is safe to play sound concurrently
type safeSoundAlert struct {
	SendSoundAlert  func()
	IsAlertNeeded   chan bool
	LastTriggerTime time.Time
}

// NRConfig addition for config required for New Relic
type NRConfig struct {
	APIKey            string                     `json:"api_key"`
	APPSets           []APPSet                   `json:"app_sets"`
	HealthSensitivity assembly.HealthSensitivity `json:"health_sensitivity"`
	EnabledSoundAlert bool                       `json:"enabled_sound_alert"`
	SurveyTime        int                        `json:"survey_time"`
}

// APPSet registered in New Relic
type APPSet struct {
	AppDetails
	Nested []APPSet `json:"nested"`
}

// AppDetails detailed information about application in New Relic
type AppDetails struct {
	Name         string   `json:"name"`
	ID           string   `json:"id"`
	RelicMetrics []string `json:"relic_metrics"`
}

// NewProcessor returns instance with implemented interface from core a IProcessor
func NewProcessor() core.IProcessor {
	return new(NewRelicProcessor)
}

// ParseConfig implementation of core.IProcessor
func (nrp *NewRelicProcessor) ParseConfig(pathToConfig string) {
	var relicConf NRConfig

	configFile, err := os.Open(pathToConfig)
	defer configFile.Close()

	if err != nil {
		log.Fatalln(err.Error())
	}

	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&relicConf)

	nrp.Config = relicConf
}

// Prepare implementation of core.IProcessor
func (nrp *NewRelicProcessor) Prepare() {
	// Prepare New Relic worker
	relicWorker, relicWorkerErr := worker.NewRelicWorker(nrp.Config.APIKey)
	if relicWorkerErr != nil {
		log.Fatalf("%s: %v", nrTag, relicWorkerErr.Error())
	}

	nrp.relicWorker = relicWorker

	// Create graph with given config structure
	configApps := make([]assembly.InfraObject, len(nrp.Config.APPSets))

	// Map to structure
	for i, set := range nrp.Config.APPSets {
		configApps[i] = assembly.InfraObject{
			Name: set.Name,
			Details: AppDetails{
				Name:         set.Name,
				ID:           set.ID,
				RelicMetrics: set.RelicMetrics,
			},
		}

		for _, nested := range set.Nested {
			nestedObj := assembly.InfraObject{
				Name:              nested.Name,
				NestedInfraObject: nil,
				Details: AppDetails{
					Name:         nested.Name,
					ID:           nested.ID,
					RelicMetrics: nested.RelicMetrics,
				},
			}

			configApps[i].NestedInfraObject = append(configApps[i].NestedInfraObject, nestedObj)
		}
	}

	// Build graph
	nrp.graph = assembly.MakeInfrastructureGraph(configApps)
	nrp.vRegionGraph = assembly.ConvertToVizceral(nrp.graph, nrp.Config.HealthSensitivity)

	// Load miscellaneous plugin to alerting by sound
	if nrp.Config.EnabledSoundAlert {
		mod, err := plugin.Open("./sound.so")
		if err != nil {
			return
		}

		SendSoundAlert, err := mod.Lookup("SendSoundAlert")
		if err != nil {
			return
		}

		log.Printf("%s: Sound plugin loaded", nrTag)

		nrp.sound.IsAlertNeeded = make(chan bool, 1)
		nrp.sound.SendSoundAlert = SendSoundAlert.(func())
	}
}

// GetLastAppGraph implementation of core.IProcessor
func (nrp *NewRelicProcessor) GetLastAppGraph() structure.VRegionGraph {
	return nrp.vRegionGraph
}

// Run implementation of core.IProcessor
func (nrp *NewRelicProcessor) Run() {
	go nrp.handleSoundAlert()

	for range time.Tick(time.Duration(nrp.Config.SurveyTime) * time.Millisecond) {
		nrp.handleMonitoring()
	}

	return
}

func (nrp *NewRelicProcessor) handleMonitoring() {
	// relicAppSets amount of applications to monitor in Ne Relic system
	infraObjects := nrp.graph.GetAllVertices()
	relicAppSets := int8(len(infraObjects))

	// Each worker will monitor one application
	var wg sync.WaitGroup
	var w int8

	for w = 0; w < relicAppSets; w++ {
		wg.Add(1)

		go func(w int8) {
			defer wg.Done()

			appVertex := nrp.graph.GetVertex(infraObjects[w])
			log.Printf("%s: Working with [%s]", nrTag, appVertex.Name)
			if appVertex.SystemDetails == nil {
				return
			}

			fetchedMetrics := nrp.relicWorker.CollectApplicationHostMetrics(
				appVertex.SystemDetails.(AppDetails).ID,
				appVertex.SystemDetails.(AppDetails).RelicMetrics,
			)

			var totalVertexRequests float32
			for _, host := range fetchedMetrics.HostMetrics {
				log.Printf("%s: Collecting metrics from host %s - %s", nrTag, appVertex.Name, host.HostName)
				uniqueHostName := host.HostName + " - " + appVertex.SystemDetails.(AppDetails).Name
				nrp.graph.AddEdge(
					graph.VertexLabel(appVertex.SystemDetails.(AppDetails).Name),
					graph.VertexLabel(uniqueHostName),
				)

				edge := nrp.graph.GetVertexEdges(graph.VertexLabel(appVertex.SystemDetails.(AppDetails).Name))
				appEdge := edge[graph.VertexLabel(uniqueHostName)]

				appEdge.Metrics.Normal += host.Metrics.Normal
				appEdge.Metrics.Warning += host.Metrics.Warning
				appEdge.Metrics.Danger += host.Metrics.Danger
				totalVertexRequests += host.Metrics.Normal + host.Metrics.Warning + host.Metrics.Danger

				appEdge.Class = assembly.GetMetricHealth(nrp.Config.HealthSensitivity, appEdge.Metrics)

				hostVertex := nrp.graph.GetVertex(graph.VertexLabel(uniqueHostName))
				hostVertex.Notices = make([]structure.VNotice, 0)
				hostVertex.Updated = time.Now().UnixNano()
				hostVertex.MaxVolume = float64(host.Metrics.Normal + host.Metrics.Warning + host.Metrics.Danger)

				if appEdge.Class == structure.VDanger {
					go func() {
						nrp.sound.IsAlertNeeded <- true
					}()
				}

				if isNeedToNotice, notice := assembly.GetHealthNotice(appEdge.Class); isNeedToNotice {
					appVertex.Notices = append(hostVertex.Notices, notice)
					hostVertex.Notices = append(hostVertex.Notices, notice)
				}
			}

			appVertex.Updated = time.Now().UnixNano()
			appVertex.MaxVolume = float64(totalVertexRequests * 5)
			appVertex.Renderer = structure.VRegionRenderer
			appVertex.Layout = structure.VLTRTreeLayout
		}(w)
	}

	wg.Wait()

	nrp.vRegionGraph = assembly.ConvertToVizceral(nrp.graph, nrp.Config.HealthSensitivity)
}

// handleSoundAlert for noise making
func (nrp *NewRelicProcessor) handleSoundAlert() {
	if nrp.sound.SendSoundAlert == nil {
		return
	}

	if <-nrp.sound.IsAlertNeeded {
		now := time.Now().UTC()
		passedTime := now.Sub(nrp.sound.LastTriggerTime)

		if passedTime >= 3 * time.Minute {
			log.Printf("%s: ----- Send sound alert -----", nrTag)

			nrp.sound.SendSoundAlert()
			nrp.sound.LastTriggerTime = now
		}
	}
}
