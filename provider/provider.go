package provider

import "github.com/LinMAD/BitAccretion/model"

// IProvider general interface to provide data for dashboard
type IProvider interface {
	// ParseConfig from file to structure for processor needs, like API keys, app ids, metrics names etc.
	ParseConfig(pathToConfig string) error
	// Prepare must setup provider before GetDispatchGraph()
	// do validating, relating or other processes for provider needs before execution
	Prepare() error

	// GetDispatchGraph executes provider to get graph with data
	GetDispatchGraph() (model.Graph, error)
	// GetHealthCheck must provide if plugin still can work
	// example API not reachable or plugin has errors and it must be restarted
	GetHealthCheck() model.HealthState
}
