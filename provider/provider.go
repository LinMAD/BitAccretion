package provider

import "github.com/LinMAD/BitAccretion/model"

// IProvider general interface to provide data for dashboard
type IProvider interface {
	// NewProvider must return implementation of interface
	NewProvider() IProvider
	// LoadConfig from file to structure for processor needs, like API keys, app ids, metrics names etc.
	LoadConfig(pathToConfig string) error
	// Boot must setup provider before DispatchGraph()
	// do validating, relating or other processes for provider needs before execution
	Boot() error

	// DispatchGraph executes provider to get graph with data
	DispatchGraph() (model.Graph, error)
	// ProvideHealth must provide if plugin still can work
	// example API not reachable or plugin has errors and it must be restarted
	ProvideHealth() model.HealthState
}
