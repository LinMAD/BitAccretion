package core

import (
	"github.com/LinMAD/BitAccretion/core/assembly/structure"
)

// IProcessor it's a basic interface for procession monitoring system for vizceral
type IProcessor interface {
	// ParseConfig from file to structure for processor needs, like API keys, app ids, metrics names etc.
	ParseConfig(pathToConfig string)
	// Prepare must setup processor before execution, validating, relating or washered processor needs before execution
	Prepare()
	// Run execute procession of monitoring
	Run()
	// GetLastAppGraph must return vizceral graph structure
	GetLastAppGraph() structure.VRegionGraph
}
