package provider

import (
	"github.com/LinMAD/BitAccretion/logger"
	"github.com/LinMAD/BitAccretion/model"
)

// IProvider general interface to provide data for dashboard
type IProvider interface {
	// LoadConfig from file to structure for processor needs, like API keys, app ids, metrics names etc.
	LoadConfig(pathToConfig string) error
	// Boot must setup provider before FetchNewData()
	// do validating, relating or other processes for provider needs before execution
	Boot() error
	// DispatchGraph must prepared graph of monitoring
	DispatchGraph() (model.Graph, error)
	// FetchNewData executes provider to get graph with data
	FetchNewData(l logger.ILogger) (model.Graph, error)
	// ProvideHealth must provide if plugin still can work
	// example API not reachable or plugin has errors and it must be restarted
	ProvideHealth() model.HealthState
}
