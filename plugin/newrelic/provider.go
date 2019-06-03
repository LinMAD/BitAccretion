package newrelic

import (
	"github.com/LinMAD/BitAccretion/model"
	"github.com/LinMAD/BitAccretion/plugin/newrelic/worker"
	"github.com/LinMAD/BitAccretion/provider"
)

// ProviderNewRelic base structure for plugin workflow
type ProviderNewRelic struct {
	Config NRConfig
	worker *worker.RelicWorker
}

// NewProvider for new relic API
func (nr *ProviderNewRelic) NewProvider() provider.IProvider {
	panic("implement me")
}

// LoadConfig of new relic plugin
func (nr *ProviderNewRelic) LoadConfig(pathToConfig string) error {
	panic("implement me")
}

// Boot setups all configuration and dependencies
func (nr *ProviderNewRelic) Boot() error {
	panic("implement me")
}

// DispatchGraph returns latest assembled graph
func (nr *ProviderNewRelic) DispatchGraph() (model.Graph, error) {
	panic("implement me")
}

// ProvideHealth will return health of plugin if New Relic API alive\reachable
func (nr *ProviderNewRelic) ProvideHealth() model.HealthState {
	panic("implement me")
}

// NewProvider returns instance with implemented interface
func NewProvider() provider.IProvider {
	return new(ProviderNewRelic)
}
