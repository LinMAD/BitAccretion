package worker

import (
	"fmt"
	"sync"

	"github.com/LinMAD/BitAccretion/extension/newrelic/client"
	"github.com/LinMAD/BitAccretion/logger"
	"github.com/LinMAD/BitAccretion/model"
)

const relicTag = "Relic worker"

// RelicWorker instance
type RelicWorker struct {
	relicClient *client.NRelicClient
}

// FetchedNewRelicData collection with the app information
type FetchedNewRelicData struct {
	HostMetrics []hostMetricsData
}

type hostMetricsData struct {
	HostName string
	HostID   int
	Metrics  model.SystemMetric
}

// NewRelicWorker initialize new worker for relic
func NewRelicWorker(APIKey string) (*RelicWorker, error) {
	rw := &RelicWorker{
		relicClient: client.NewRelicClient(APIKey),
	}

	isAuth, err := rw.relicClient.Authenticate()
	if err != nil || !isAuth {
		return nil, fmt.Errorf("%s: Unable to pass authetication with key: %s", relicTag, APIKey)
	}

	return rw, nil
}

// CollectApplicationHostMetrics returns all collected metrics for application
func (w *RelicWorker) CollectApplicationHostMetrics(log logger.ILogger, appID string, metrics []string) *FetchedNewRelicData {
	fetchedData := &FetchedNewRelicData{}
	fetchedData.HostMetrics = make([]hostMetricsData, 0)

	hosts := w.relicClient.GetApplicationHost(appID)
	hLen := len(hosts.AppsHosts)

	// TODO Some where here could be dead lock, mb make timout with context cancel of all go threads
	var wg sync.WaitGroup
	for group := 0; group < hLen; group++ {
		wg.Add(1)

		go func(group int) {
			defer wg.Done()

			host := hosts.AppsHosts[group]

			hostMetrics := w.relicClient.GetHostMetricData(appID, host.HostID, metrics)
			metricsData := model.SystemMetric{}

			for _, metrics := range hostMetrics.Data.Metrics {
				for _, rates := range metrics.Timeslices {
					metricsData.RequestCount += rates.Values.Requests
					metricsData.ErrorCount += rates.Values.Errors
				}
			}

			hostMetricsData := hostMetricsData{}
			hostMetricsData.HostName = host.Host
			hostMetricsData.HostID = host.HostID
			hostMetricsData.Metrics = metricsData

			fetchedData.HostMetrics = append(fetchedData.HostMetrics, hostMetricsData)
		}(group)
	}

	wg.Wait()
	log.Debug(fmt.Sprintf("Traversed NewRelic APP ID: %s and all instances host (%d)", appID, hLen))

	return fetchedData
}
