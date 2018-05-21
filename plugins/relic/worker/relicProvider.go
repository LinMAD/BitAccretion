package worker

import (
	"fmt"
	"github.com/LinMAD/BitAccretion/core/assembly/structure"
	"github.com/LinMAD/BitAccretion/plugins/relic/client"
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
	Metrics  structure.VMetricLevels
}

// NewRelicWorker initialize new worker for relic
func NewRelicWorker(APIKey string) (*RelicWorker, error) {
	rw := &RelicWorker{
		relicClient: client.NewRelicClient(APIKey),
	}

	if !rw.relicClient.Authenticate() {
		return nil, fmt.Errorf("%s: Unable to pass authetication with key: %s", relicTag, APIKey)
	}

	return rw, nil
}

// CollectApplicationHostMetrics returns all collected metrics for application
func (w *RelicWorker) CollectApplicationHostMetrics(appID string, metrics []string) *FetchedNewRelicData {
	fetchedData := &FetchedNewRelicData{}
	fetchedData.HostMetrics = make([]hostMetricsData, 0)

	hosts := w.relicClient.GetApplicationHost(appID)
	for _, host := range hosts.AppsHosts {
		hostMetrics := w.relicClient.GetHostMetricData(appID, host.HostID, metrics)
		metricsData := structure.VMetricLevels{}

		for _, metrics := range hostMetrics.Data.Metrics {
			for _, rates := range metrics.Timeslices {
				metricsData.Normal += rates.Values.Requests
				metricsData.Danger += rates.Values.Errors
			}
		}

		hostMetricsData := hostMetricsData{}
		hostMetricsData.HostName = host.Host
		hostMetricsData.HostID = host.HostID
		hostMetricsData.Metrics = metricsData

		fetchedData.HostMetrics = append(fetchedData.HostMetrics, hostMetricsData)
	}

	return fetchedData
}
